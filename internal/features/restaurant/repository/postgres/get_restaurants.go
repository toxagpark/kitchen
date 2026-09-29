package restaurant_pg_repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	restaurant_model "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/restaurant"
)

func (r *Repository) GetRestaurants(
	c *gin.Context,
	latitude float64,
	longitude float64,
	page int,
	pageSize int,
) ([]restaurant_model.RestaurantWithDistance, error) {
	const query = `
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			delivery_radius_m,
			opening_time,
			closing_time,
			6371000 * 2 * ASIN(
				SQRT(
					POWER(SIN(RADIANS(latitude - $1) / 2), 2) +
					COS(RADIANS($1)) *
					COS(RADIANS(latitude)) *
					POWER(SIN(RADIANS(longitude - $2) / 2), 2)
				)
			) AS distance_m
		FROM kitchen.restaurants
		ORDER BY distance_m
		LIMIT $3
		OFFSET $4
	`

	offset := page * pageSize

	rows, err := r.db.Query(
		c.Request.Context(),
		query,
		latitude,
		longitude,
		pageSize,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get restaurants: %w", err)
	}
	defer rows.Close()

	result := make([]restaurant_model.RestaurantWithDistance, 0)

	for rows.Next() {
		var restaurant domain.Restaurant
		var distanceM float64

		if err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.Latitude,
			&restaurant.Longitude,
			&restaurant.DeliveryRadiusM,
			&restaurant.OpeningTime,
			&restaurant.ClosingTime,
			&distanceM,
		); err != nil {
			return nil, fmt.Errorf("scan restaurant: %w", err)
		}

		result = append(result, restaurant_model.RestaurantWithDistance{
			Restaurant: restaurant,
			DistanceM:  int(distanceM),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate restaurants: %w", err)
	}

	return result, nil
}
