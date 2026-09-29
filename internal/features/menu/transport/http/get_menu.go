package menu_transport_http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	core_errors "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/errors"
	core_logger "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/logger"
	core_http_response "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/transport/http/response"
)

func (h *MenuHTTPHandler) GetMenu(c *gin.Context) {
	log := core_logger.FromContext(c.Request.Context())
	rw := core_http_response.NewHTTPResponseHandler(log, c)

	restaurantID, err := uuid.Parse(c.Param("restaurant_id"))
	if err != nil {
		rw.ErrorResponse(
			fmt.Errorf("invalid restaurant_id: %w", core_errors.ErrInvalidArgument),
			"parse restaurant_id",
		)
		return
	}

	items, err := h.menuService.GetMenu(c, restaurantID)
	if err != nil {
		rw.ErrorResponse(err, "failed to get menu")
		return
	}

	response := make([]MenuItemResponse, 0, len(items))

	for _, item := range items {
		response = append(response, MenuItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			IsAvailable: item.IsAvailable,
		})
	}

	c.JSON(http.StatusOK, response)
}

type MenuItemResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	IsAvailable bool      `json:"is_available"`
}
