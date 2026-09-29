package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID      uuid.UUID
	Name    string
	Address string

	Latitude        float64
	Longitude       float64
	DeliveryRadiusM int

	OpeningTime time.Time
	ClosingTime time.Time
}

type RestaurantWithContext struct {
	Restaurant        Restaurant
	DistanceM         int
	IsOpen            bool
	DeliveryAvailable bool
}

func NewRestaurant(
	id uuid.UUID,
	name string,
	address string,
	latitude float64,
	longitude float64,
	deliveryRadiusM int,
	openingTime time.Time,
	closingTime time.Time,
) Restaurant {
	return Restaurant{
		ID:              id,
		Name:            name,
		Address:         address,
		Latitude:        latitude,
		Longitude:       longitude,
		DeliveryRadiusM: deliveryRadiusM,
		OpeningTime:     openingTime,
		ClosingTime:     closingTime,
	}
}

func (r Restaurant) Validate() error {
	if r.ID == uuid.Nil {
		return fmt.Errorf("invalid restaurant id: %w", ErrInvalidRestaurant)
	}

	if r.Name == "" {
		return fmt.Errorf("invalid restaurant name: %w", ErrInvalidRestaurant)
	}

	if r.Address == "" {
		return fmt.Errorf("invalid restaurant address: %w", ErrInvalidRestaurant)
	}

	if r.Latitude < -90 || r.Latitude > 90 {
		return fmt.Errorf("invalid latitude: %w", ErrInvalidRestaurant)
	}

	if r.Longitude < -180 || r.Longitude > 180 {
		return fmt.Errorf("invalid longitude: %w", ErrInvalidRestaurant)
	}

	if r.DeliveryRadiusM <= 0 {
		return fmt.Errorf("invalid delivery radius: %w", ErrInvalidRestaurant)
	}

	return nil
}

func (r Restaurant) IsOpen() bool {
	now := time.Now()

	currentMinutes := now.Hour()*60 + now.Minute()
	openingMinutes := r.OpeningTime.Hour()*60 + r.OpeningTime.Minute()
	closingMinutes := r.ClosingTime.Hour()*60 + r.ClosingTime.Minute()

	if openingMinutes < closingMinutes {
		return currentMinutes >= openingMinutes &&
			currentMinutes < closingMinutes
	}

	return currentMinutes >= openingMinutes ||
		currentMinutes < closingMinutes
}
