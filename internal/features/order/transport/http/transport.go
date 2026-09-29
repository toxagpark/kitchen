package order_transport_http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order"
)

type OrderHTTPHandler struct {
	orderService OrderService
}

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		restaurantID uuid.UUID,
		requestedItems []order.CreateOrderItem,
	) (domain.Order, error)
}

func NewOrderHTTPHandler(
	orderService OrderService,
) *OrderHTTPHandler {
	return &OrderHTTPHandler{
		orderService: orderService,
	}
}

func (h *OrderHTTPHandler) Routes(rg *gin.RouterGroup) {
	rg.POST("/orders", h.CreateOrder)
}
