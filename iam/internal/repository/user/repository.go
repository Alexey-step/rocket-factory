package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	iamRepo "github.com/Alexey-step/rocket-factory/iam/internal/repository"
)

var _ iamRepo.UserRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}
