package menu_transport_http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
)

type MenuHTTPHandler struct {
	menuService MenuService
}

type MenuService interface {
	GetMenu(
		c *gin.Context,
		restaurantID uuid.UUID,
	) ([]domain.MenuItem, error)
}

func NewRestaurantHTTPHandler(
	menuService MenuService,
) *MenuHTTPHandler {
	return &MenuHTTPHandler{
		menuService: menuService,
	}

}

func (h *MenuHTTPHandler) Routes(rg *gin.RouterGroup) {
	rg.GET("/restaurants/:restaurant_id/menu", h.GetMenu)
}
