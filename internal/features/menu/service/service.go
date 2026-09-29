package menu_service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/domain"
)

type Service struct {
	repository MenuRepository
}

type MenuRepository interface {
	GetMenu(
		c *gin.Context,
		restaurantID uuid.UUID,
	) ([]domain.MenuItem, error)
}

func NewService(repository MenuRepository) *Service {
	return &Service{
		repository: repository,
	}
}
