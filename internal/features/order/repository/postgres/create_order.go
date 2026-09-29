package order_pg_repository

import (
	"context"
	"fmt"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	core_postgres_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool"
)

func (r *Repository) CreateOrder(
	ctx context.Context,
	tx core_postgres_pool.Tx,
	order domain.Order,
) error {
	const query = `
		INSERT INTO kitchen.orders (
			id,
			restaurant_id,
			status,
			total_price
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(
		ctx,
		query,
		order.ID,
		order.RestaurantID,
		order.Status,
		order.TotalPrice,
	)
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
}
