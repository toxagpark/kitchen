package restaurants_transport_http

import (
	"github.com/gin-gonic/gin"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
)

type RestaurantsHTTPHandler struct {
	restaurantsService RestaurantsService
}

type RestaurantsService interface {
	GetRestaurants(
		c *gin.Context,
		latitude float64,
		longitude float64,
		page int,
		pageSize int,
	) ([]domain.RestaurantWithContext, error)
}

func NewRestaurantHTTPHandler(
	restaurantService RestaurantsService,
) *RestaurantsHTTPHandler {
	return &RestaurantsHTTPHandler{
		restaurantsService: restaurantService,
	}

}

func (h *RestaurantsHTTPHandler) Routes(rg *gin.RouterGroup) {
	rg.GET("/restaurants", h.GetRestaurants)
}
