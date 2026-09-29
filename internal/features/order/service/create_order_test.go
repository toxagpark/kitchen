package order_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	core_errors "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order/service/mocks"
)

func TestCreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	txManager := mocks.NewMockTransactionManager(ctrl)
	restaurantClient := mocks.NewMockRestaurantClient(ctrl)
	tx := mocks.NewMockTx(ctrl)

	service := NewService(
		txManager,
		repository,
		restaurantClient,
	)

	ctx := context.Background()

	restaurantID := uuid.New()
	menuItemID := uuid.New()

	requestedItems := []order.CreateOrderItem{
		{
			MenuItemID: menuItemID,
			Quantity:   2,
		},
	}

	menuItems := []domain.MenuItem{
		{
			ID:           menuItemID,
			RestaurantID: restaurantID,
			Name:         "Пицца",
			Description:  "Маргарита",
			Price:        50000,
			IsAvailable:  true,
		},
	}

	repository.EXPECT().
		GetMenuItems(ctx, restaurantID, []uuid.UUID{menuItemID}).
		Return(menuItems, nil)

	txManager.EXPECT().
		Begin(ctx).
		Return(tx, nil)

	repository.EXPECT().
		CreateOrder(ctx, tx, gomock.Any()).
		Return(nil)

	repository.EXPECT().
		CreateOrderItems(ctx, tx, gomock.Any()).
		Return(nil)

	tx.EXPECT().
		Commit(ctx).
		Return(nil)

	restaurantClient.EXPECT().
		CreateOrder(ctx, gomock.Any(), gomock.Any()).
		Return(nil)

	repository.EXPECT().
		UpdateOrderStatus(
			ctx,
			gomock.Any(),
			domain.OrderStatusAccepted,
		).
		Return(nil)

	tx.EXPECT().
		Rollback(ctx).
		Return(nil)

	result, err := service.CreateOrder(
		ctx,
		restaurantID,
		requestedItems,
	)

	require.NoError(t, err)
	require.Equal(t, restaurantID, result.RestaurantID)
	require.Equal(t, domain.OrderStatusAccepted, result.Status)
	require.Equal(t, int64(100000), result.TotalPrice)
	require.False(t, result.ID == uuid.Nil)
	require.False(t, result.CreatedAt.IsZero())
}

func TestCreateOrder_InvalidQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	txManager := mocks.NewMockTransactionManager(ctrl)
	restaurantClient := mocks.NewMockRestaurantClient(ctrl)

	service := NewService(
		txManager,
		repository,
		restaurantClient,
	)

	ctx := context.Background()
	restaurantID := uuid.New()
	menuItemID := uuid.New()

	requestedItems := []order.CreateOrderItem{
		{
			MenuItemID: menuItemID,
			Quantity:   0,
		},
	}

	result, err := service.CreateOrder(
		ctx,
		restaurantID,
		requestedItems,
	)

	require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	require.Equal(t, domain.Order{}, result)
}

func TestCreateOrder_MenuItemUnavailable(t *testing.T) {
	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	txManager := mocks.NewMockTransactionManager(ctrl)
	restaurantClient := mocks.NewMockRestaurantClient(ctrl)

	service := NewService(
		txManager,
		repository,
		restaurantClient,
	)

	ctx := context.Background()

	restaurantID := uuid.New()
	menuItemID := uuid.New()

	requestedItems := []order.CreateOrderItem{
		{
			MenuItemID: menuItemID,
			Quantity:   1,
		},
	}

	menuItems := []domain.MenuItem{
		{
			ID:           menuItemID,
			RestaurantID: restaurantID,
			Name:         "Пицца",
			Description:  "Маргарита",
			Price:        50000,
			IsAvailable:  false,
		},
	}

	repository.EXPECT().
		GetMenuItems(ctx, restaurantID, []uuid.UUID{menuItemID}).
		Return(menuItems, nil)

	result, err := service.CreateOrder(
		ctx,
		restaurantID,
		requestedItems,
	)

	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Equal(t, domain.Order{}, result)
}

func TestCreateOrder_RestaurantRejects(t *testing.T) {
	ctrl := gomock.NewController(t)

	repository := mocks.NewMockRepository(ctrl)
	txManager := mocks.NewMockTransactionManager(ctrl)
	restaurantClient := mocks.NewMockRestaurantClient(ctrl)
	tx := mocks.NewMockTx(ctrl)

	service := NewService(
		txManager,
		repository,
		restaurantClient,
	)

	ctx := context.Background()

	restaurantID := uuid.New()
	menuItemID := uuid.New()

	requestedItems := []order.CreateOrderItem{
		{
			MenuItemID: menuItemID,
			Quantity:   2,
		},
	}

	menuItems := []domain.MenuItem{
		{
			ID:           menuItemID,
			RestaurantID: restaurantID,
			Name:         "Пицца",
			Description:  "Маргарита",
			Price:        50000,
			IsAvailable:  true,
		},
	}

	repository.EXPECT().
		GetMenuItems(ctx, restaurantID, []uuid.UUID{menuItemID}).
		Return(menuItems, nil)

	txManager.EXPECT().
		Begin(ctx).
		Return(tx, nil)

	repository.EXPECT().
		CreateOrder(ctx, tx, gomock.Any()).
		Return(nil)

	repository.EXPECT().
		CreateOrderItems(ctx, tx, gomock.Any()).
		Return(nil)

	tx.EXPECT().
		Commit(ctx).
		Return(nil)

	restaurantClient.EXPECT().
		CreateOrder(ctx, gomock.Any(), gomock.Any()).
		Return(errors.New("restaurant unavailable"))

	repository.EXPECT().
		UpdateOrderStatus(
			ctx,
			gomock.Any(),
			domain.OrderStatusRejected,
		).
		Return(nil)

	tx.EXPECT().
		Rollback(ctx).
		Return(nil)

	result, err := service.CreateOrder(
		ctx,
		restaurantID,
		requestedItems,
	)

	require.NoError(t, err)
	require.Equal(t, domain.OrderStatusRejected, result.Status)
	require.Equal(t, int64(100000), result.TotalPrice)
	require.Equal(t, restaurantID, result.RestaurantID)
}
