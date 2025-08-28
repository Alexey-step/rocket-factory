package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/order/internal/config"
	customMiddleware "github.com/Alexey-step/rocket-factory/order/internal/middleware"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	orderV1 "github.com/Alexey-step/rocket-factory/shared/pkg/openapi/order/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	// HTTP сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- errors.Errorf("http server crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
		a.initMigrations,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initMigrations(ctx context.Context) error {
	err := a.diContainer.Migrator(ctx).Up(ctx)
	if err != nil {
		logger.Error(ctx, "ошибка миграции базы данных", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	serv, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		panic(fmt.Sprintf("ошибка создания сервера OpenAPI: %v\n", err))
	}

	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.RequestLogger)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", serv)

	a.httpServer = &http.Server{
		Addr:        config.AppConfig().OrderHTTP.Address(),
		Handler:     r,
		ReadTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		err = a.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer running")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
