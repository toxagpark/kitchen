package restaurants_service

import (
	"github.com/gin-gonic/gin"
	restaurant_model "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/restaurant"
)

type Service struct {
	repository RestaurantRepository
}

type RestaurantRepository interface {
	GetRestaurants(c *gin.Context,
		latitude float64,
		longitude float64,
		page int,
		pageSize int,
	) ([]restaurant_model.RestaurantWithDistance, error)
}

func NewService(repo RestaurantRepository) *Service {
	return &Service{
		repository: repo,
	}
}
