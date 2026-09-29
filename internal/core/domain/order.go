package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusAccepted OrderStatus = "accepted"
	OrderStatusRejected OrderStatus = "rejected"
)

type Order struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Status       OrderStatus
	TotalPrice   int64
	CreatedAt    time.Time
}

func NewOrder(
	id uuid.UUID,
	restaurantID uuid.UUID,
	status OrderStatus,
	totalPrice int64,
	createdAt time.Time,
) Order {
	return Order{
		ID:           id,
		RestaurantID: restaurantID,
		Status:       status,
		TotalPrice:   totalPrice,
		CreatedAt:    createdAt,
	}
}

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	MenuItemID uuid.UUID

	Name       string
	Quantity   int
	UnitPrice  int64
	TotalPrice int64
}

func NewOrderItem(
	id uuid.UUID,
	orderID uuid.UUID,
	menuItemID uuid.UUID,

	name string,
	quantity int,
	unitPrice int64,
	totalPrice int64,
) OrderItem {
	return OrderItem{
		ID:         id,
		OrderID:    orderID,
		MenuItemID: menuItemID,
		Name:       name,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}
}
