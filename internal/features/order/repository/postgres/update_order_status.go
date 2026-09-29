package order_pg_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
)

func (r *Repository) UpdateOrderStatus(
	ctx context.Context,
	orderID uuid.UUID,
	status domain.OrderStatus,
) error {
	const query = `
		UPDATE kitchen.orders
		SET status = $1
		WHERE id = $2
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		status,
		orderID,
	)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}
