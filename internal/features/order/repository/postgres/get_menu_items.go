package order_pg_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
)

func (r *Repository) GetMenuItems(
	ctx context.Context,
	restaurantID uuid.UUID,
	menuItemIDs []uuid.UUID,
) ([]domain.MenuItem, error) {
	const query = `
		SELECT
			id,
			restaurant_id,
			name,
			description,
			price,
			is_available
		FROM kitchen.menu_items
		WHERE restaurant_id = $1
		  AND id = ANY($2)
		ORDER BY name
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		restaurantID,
		menuItemIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("get menu items: %w", err)
	}
	defer rows.Close()

	result := make([]domain.MenuItem, 0, len(menuItemIDs))

	for rows.Next() {
		var item domain.MenuItem

		if err := rows.Scan(
			&item.ID,
			&item.RestaurantID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.IsAvailable,
		); err != nil {
			return nil, fmt.Errorf("scan menu item: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate menu items: %w", err)
	}

	return result, nil
}
