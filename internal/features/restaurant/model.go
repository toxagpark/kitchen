package restaurant_model

import "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"

type RestaurantWithDistance struct {
	Restaurant domain.Restaurant
	DistanceM  int
}
