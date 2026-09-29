package order

import "github.com/google/uuid"

type CreateOrderItem struct {
	MenuItemID uuid.UUID
	Quantity   int
}
