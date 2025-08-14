package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/Alexey-step/rocket-factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}
