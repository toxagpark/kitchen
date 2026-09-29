package order_transport_http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	core_errors "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/errors"
	core_logger "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/logger"
	core_http_response "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/transport/http/response"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order"
)

func (h *OrderHTTPHandler) CreateOrder(c *gin.Context) {
	log := core_logger.FromContext(c.Request.Context())
	rw := core_http_response.NewHTTPResponseHandler(log, c)

	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		rw.ErrorResponse(core_errors.ErrInvalidArgument, "error bind json")
		return
	}

	requestedItems := make(
		[]order.CreateOrderItem,
		0,
		len(req.Items),
	)

	for _, item := range req.Items {
		requestedItems = append(
			requestedItems,
			order.CreateOrderItem{
				MenuItemID: item.MenuItemID,
				Quantity:   item.Quantity,
			},
		)
	}

	order, err := h.orderService.CreateOrder(
		c,
		req.RestaurantID,
		requestedItems,
	)
	if err != nil {
		rw.ErrorResponse(err, "failed to create order")
		return
	}

	c.JSON(http.StatusCreated, OrderResponse{
		ID:           order.ID,
		RestaurantID: order.RestaurantID,
		Status:       string(order.Status),
		TotalPrice:   order.TotalPrice,
		CreatedAt:    order.CreatedAt,
	})
}

type CreateOrderRequest struct {
	RestaurantID uuid.UUID                `json:"restaurant_id" binding:"required"`
	Items        []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
}

type CreateOrderItemRequest struct {
	MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
	Quantity   int       `json:"quantity" binding:"required,min=1"`
}

type OrderResponse struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Status       string    `json:"status"`
	TotalPrice   int64     `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
}
