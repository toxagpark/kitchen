package order_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
	core_postgres_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool"
)

type Service struct {
	tx               TransactionManager
	repository       Repository
	restaurantClient RestaurantClient
}

func NewService(
	tx TransactionManager,
	repository Repository,
	restaurantClient RestaurantClient,
) *Service {
	return &Service{
		tx:               tx,
		repository:       repository,
		restaurantClient: restaurantClient,
	}
}

type Repository interface {
	GetMenuItems(
		ctx context.Context,
		restaurantID uuid.UUID,
		menuItemIDs []uuid.UUID,
	) ([]domain.MenuItem, error)

	CreateOrder(
		ctx context.Context,
		tx core_postgres_pool.Tx,
		order domain.Order,
	) error

	CreateOrderItems(
		ctx context.Context,
		tx core_postgres_pool.Tx,
		items []domain.OrderItem,
	) error

	UpdateOrderStatus(
		ctx context.Context,
		orderID uuid.UUID,
		status domain.OrderStatus,
	) error
}

type TransactionManager interface {
	Begin(ctx context.Context) (core_postgres_pool.Tx, error)
}

type RestaurantClient interface {
	CreateOrder(
		ctx context.Context,
		order domain.Order,
		items []domain.OrderItem,
	) error
}
