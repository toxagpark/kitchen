package restaurants_service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	core_errors "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/errors"
)

const (
	minLatitude  = -90.0
	maxLatitude  = 90.0
	minLongitude = -180.0
	maxLongitude = 180.0
)

func (s *Service) GetRestaurants(
	c *gin.Context,
	latitude float64,
	longitude float64,
	page int,
	pageSize int,
) ([]domain.RestaurantWithContext, error) {
	if pageSize == 0 {
		pageSize = 20
	}

	if err := validateCoordinates(latitude, longitude); err != nil {
		return nil, err
	}

	restaurants, err := s.repository.GetRestaurants(
		c,
		latitude,
		longitude,
		page,
		pageSize,
	)
	if err != nil {
		return nil, err
	}

	result := make([]domain.RestaurantWithContext, 0, len(restaurants))

	for _, restaurant := range restaurants {
		if err := restaurant.Restaurant.Validate(); err != nil {
			return nil, err
		}
		result = append(result, domain.RestaurantWithContext{
			Restaurant:        restaurant.Restaurant,
			DistanceM:         restaurant.DistanceM,
			IsOpen:            restaurant.Restaurant.IsOpen(),
			DeliveryAvailable: restaurant.DistanceM <= restaurant.Restaurant.DeliveryRadiusM,
		})
	}

	return result, nil
}

func validateCoordinates(latitude, longitude float64) error {
	if latitude < minLatitude || latitude > maxLatitude {
		return fmt.Errorf("invalid latitude: %w", core_errors.ErrInvalidArgument)
	}

	if longitude < minLongitude || longitude > maxLongitude {
		return fmt.Errorf("invalid longitude: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
