package domain

import "github.com/google/uuid"

type MenuItem struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Name         string
	Description  string
	Price        int64
	IsAvailable  bool
}
