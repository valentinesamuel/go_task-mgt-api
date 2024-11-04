package user

import (
	"context"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, User *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, User *models.User) (*models.User, error)
	Delete(ctx context.Context, id uint) (*models.User, error)
}
