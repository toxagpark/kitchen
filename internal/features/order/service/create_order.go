package order_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	core_errors "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/errors"
	core_postgres_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order"
)

func (s *Service) CreateOrder(
	ctx context.Context,
	restaurantID uuid.UUID,
	requestedItems []order.CreateOrderItem,
) (domain.Order, error) {
	if err := validateRequestedItems(requestedItems); err != nil {
		return domain.Order{}, fmt.Errorf("validate order items: %w", err)
	}

	menuItemIDs := getMenuItemIDs(requestedItems)

	menuItems, err := s.repository.GetMenuItems(
		ctx,
		restaurantID,
		menuItemIDs,
	)
	if err != nil {
		return domain.Order{}, err
	}

	if len(menuItems) != len(requestedItems) {
		return domain.Order{}, core_errors.ErrNotFound
	}

	menuItemsByID := mapMenuItemsByID(menuItems)

	orderID := uuid.New()

	orderItems, totalPrice, err := buildOrderItems(
		orderID,
		requestedItems,
		menuItemsByID,
	)
	if err != nil {
		return domain.Order{}, err
	}

	order := domain.Order{
		ID:           orderID,
		RestaurantID: restaurantID,
		Status:       domain.OrderStatusPending,
		TotalPrice:   totalPrice,
		CreatedAt:    time.Now(),
	}

	tx, err := s.tx.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	return s.createOrder(
		ctx,
		tx,
		order,
		orderItems,
	)
}

func validateRequestedItems(
	items []order.CreateOrderItem,
) error {
	if len(items) == 0 {
		return core_errors.ErrInvalidArgument
	}

	seen := make(map[uuid.UUID]struct{}, len(items))

	for _, item := range items {
		if item.Quantity <= 0 {
			return core_errors.ErrInvalidArgument
		}

		if _, exists := seen[item.MenuItemID]; exists {
			return core_errors.ErrInvalidArgument
		}

		seen[item.MenuItemID] = struct{}{}
	}

	return nil
}

func getMenuItemIDs(
	items []order.CreateOrderItem,
) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(items))

	for _, item := range items {
		result = append(result, item.MenuItemID)
	}

	return result
}

func mapMenuItemsByID(
	menuItems []domain.MenuItem,
) map[uuid.UUID]domain.MenuItem {
	result := make(map[uuid.UUID]domain.MenuItem, len(menuItems))

	for _, menuItem := range menuItems {
		result[menuItem.ID] = menuItem
	}

	return result
}

func buildOrderItems(
	orderID uuid.UUID,
	requestedItems []order.CreateOrderItem,
	menuItemsByID map[uuid.UUID]domain.MenuItem,
) ([]domain.OrderItem, int64, error) {
	orderItems := make([]domain.OrderItem, 0, len(requestedItems))

	var totalPrice int64

	for _, requestedItem := range requestedItems {
		menuItem, exists := menuItemsByID[requestedItem.MenuItemID]
		if !exists {
			return nil, 0, core_errors.ErrNotFound
		}

		if !menuItem.IsAvailable {
			return nil, 0, core_errors.ErrConflict
		}

		itemTotal := menuItem.Price * int64(requestedItem.Quantity)

		orderItems = append(orderItems, domain.OrderItem{
			ID:         uuid.New(),
			OrderID:    orderID,
			MenuItemID: menuItem.ID,
			Name:       menuItem.Name,
			Quantity:   requestedItem.Quantity,
			UnitPrice:  menuItem.Price,
			TotalPrice: itemTotal,
		})

		totalPrice += itemTotal
	}

	return orderItems, totalPrice, nil
}

func (s *Service) createOrder(
	ctx context.Context,
	tx core_postgres_pool.Tx,
	order domain.Order,
	orderItems []domain.OrderItem,
) (domain.Order, error) {
	if err := s.repository.CreateOrder(
		ctx,
		tx,
		order,
	); err != nil {
		return domain.Order{}, err
	}

	if err := s.repository.CreateOrderItems(
		ctx,
		tx,
		orderItems,
	); err != nil {
		return domain.Order{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}

	if err := s.restaurantClient.CreateOrder(
		ctx,
		order,
		orderItems,
	); err != nil {
		if updateErr := s.repository.UpdateOrderStatus(
			ctx,
			order.ID,
			domain.OrderStatusRejected,
		); updateErr != nil {
			return domain.Order{}, updateErr
		}

		order.Status = domain.OrderStatusRejected

		return order, nil
	}

	if err := s.repository.UpdateOrderStatus(
		ctx,
		order.ID,
		domain.OrderStatusAccepted,
	); err != nil {
		return domain.Order{}, err
	}

	order.Status = domain.OrderStatusAccepted

	return order, nil
}
