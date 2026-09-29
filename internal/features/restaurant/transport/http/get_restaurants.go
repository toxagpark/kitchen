package restaurants_transport_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	core_errors "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/errors"
	core_logger "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/logger"
	core_http_response "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/transport/http/response"
)

type RestaurantsRequest struct {
	Latitude  *float64 `form:"latitude" binding:"required,min=-90,max=90" example:"20"` //указатель для нуля(что б прошел)
	Longitude *float64 `form:"longitude" binding:"required,min=-180,max=180" example:"30"`
	Page      int      `form:"page" binding:"omitempty,min=0" example:"1"`
	PageSize  int      `form:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
}

type RestaurantResponse struct {
	ID                uuid.UUID `json:"id" example:"8f3a1c2e-7b4d-4f91-9c62-123456789abc"`
	Name              string    `json:"name" example:"Пицца Хаус"`
	Address           string    `json:"address" example:"ул. Ленина, 10"`
	DistanceM         int       `json:"distance_m" example:"850"`
	IsOpen            bool      `json:"is_open" example:"true"`
	DeliveryAvailable bool      `json:"delivery_available" example:"true"`
}

func newRestaurantResponse(
	id uuid.UUID,
	name string,
	address string,
	distanceM int,
	isOpen bool,
	deliveryAvailable bool,
) RestaurantResponse {
	return RestaurantResponse{
		ID:                id,
		Name:              name,
		Address:           address,
		DistanceM:         distanceM,
		IsOpen:            isOpen,
		DeliveryAvailable: deliveryAvailable,
	}
}

func restaurantDomainToResponse(rest domain.RestaurantWithContext) RestaurantResponse {
	return newRestaurantResponse(
		rest.Restaurant.ID,
		rest.Restaurant.Name,
		rest.Restaurant.Address,
		rest.DistanceM,
		rest.IsOpen,
		rest.DeliveryAvailable,
	)
}

func restaurantDomainsToResponse(rests []domain.RestaurantWithContext) []RestaurantResponse {
	resp := make([]RestaurantResponse, 0, len(rests))
	for i := 0; i < len(rests); i++ {
		rest := restaurantDomainToResponse(rests[i])
		resp = append(resp, rest)
	}
	return resp
}

func (h *RestaurantsHTTPHandler) GetRestaurants(c *gin.Context) {
	log := core_logger.FromContext(c.Request.Context())
	rw := core_http_response.NewHTTPResponseHandler(log, c)
	var req RestaurantsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		rw.ErrorResponse(core_errors.ErrInvalidArgument, "invalid latitude/longitude")
		return
	}

	restaurantWithContextDomains, err := h.restaurantsService.GetRestaurants(
		c,
		*req.Latitude,
		*req.Longitude,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		rw.ErrorResponse(err, "failed to get restaurants")
		return
	}

	resp := restaurantDomainsToResponse(restaurantWithContextDomains)

	c.JSON(http.StatusOK, resp)
}
