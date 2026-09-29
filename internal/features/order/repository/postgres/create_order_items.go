package order_pg_repository

import (
	"context"
	"fmt"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	core_postgres_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool"
)

func (r *Repository) CreateOrderItems(
	ctx context.Context,
	tx core_postgres_pool.Tx,
	items []domain.OrderItem,
) error {
	const query = `
		INSERT INTO kitchen.order_items (
			id,
			order_id,
			menu_item_id,
			name,
			quantity,
			unit_price,
			total_price
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, item := range items {
		_, err := tx.Exec(
			ctx,
			query,
			item.ID,
			item.OrderID,
			item.MenuItemID,
			item.Name,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
		)
		if err != nil {
			return fmt.Errorf("create order item: %w", err)
		}
	}

	return nil
}
