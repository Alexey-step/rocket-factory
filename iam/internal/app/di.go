package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	authV1API "github.com/Alexey-step/rocket-factory/iam/internal/api/auth/v1"
	userV1API "github.com/Alexey-step/rocket-factory/iam/internal/api/user/v1"
	"github.com/Alexey-step/rocket-factory/iam/internal/config"
	"github.com/Alexey-step/rocket-factory/iam/internal/repository"
	sessionRepository "github.com/Alexey-step/rocket-factory/iam/internal/repository/session"
	userRepository "github.com/Alexey-step/rocket-factory/iam/internal/repository/user"
	"github.com/Alexey-step/rocket-factory/iam/internal/service"
	authService "github.com/Alexey-step/rocket-factory/iam/internal/service/auth"
	userService "github.com/Alexey-step/rocket-factory/iam/internal/service/user"
	"github.com/Alexey-step/rocket-factory/platform/pkg/cache"
	"github.com/Alexey-step/rocket-factory/platform/pkg/cache/redis"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	"github.com/Alexey-step/rocket-factory/platform/pkg/migrator"
	pgMigrator "github.com/Alexey-step/rocket-factory/platform/pkg/migrator/pg"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	userV1API userV1.UserServiceServer
	authV1Api authV1.AuthServiceServer

	userService service.UserService
	authService service.AuthService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	postgresDB *pgxpool.Pool
	migrator   migrator.Migrator

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1Api == nil {
		d.authV1Api = authV1API.NewAPI(d.AuthService(ctx))
	}

	return d.authV1Api
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.SessionRepository(ctx),
			d.UserRepository(ctx),
			config.AppConfig().Session.TTL(),
		)
	}

	return d.authService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PostgresDB(ctx))
	}

	return d.userRepository
}

func (d *diContainer) SessionRepository(_ context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient())
	}

	return d.sessionRepository
}

func (d *diContainer) PostgresDB(ctx context.Context) *pgxpool.Pool {
	if d.postgresDB == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Postgres: %s\n", err.Error()))
		}

		// Проверяем соединение с базой данных
		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping Postgres: %s\n", err.Error()))
		}

		closer.AddNamed("PostgresDB", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.postgresDB = pool
	}

	return d.postgresDB
}

func (d *diContainer) Migrator(_ context.Context) migrator.Migrator {
	if d.migrator == nil {
		cfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to parse Postgres config: %s\n", err.Error()))
		}

		migrationsDir := config.AppConfig().IamGRPC.MigrationsDir()
		d.migrator = pgMigrator.NewMigrator(stdlib.OpenDB(*cfg.ConnConfig), migrationsDir)
	}

	return d.migrator
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIDLE(),
			IdleTimeout: config.AppConfig().Redis.IDLETimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}
