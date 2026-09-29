package menu_pg_repository

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
)

func (r *Repository) GetMenu(
	c *gin.Context,
	restaurantID uuid.UUID,
) ([]domain.MenuItem, error) {
	query := `
        SELECT
            id,
            restaurant_id,
            name,
            description,
            price,
            is_available
        FROM kitchen.menu_items
        WHERE restaurant_id = $1
        ORDER BY name
    `

	rows, err := r.db.Query(c, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.MenuItem, 0)

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
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}
