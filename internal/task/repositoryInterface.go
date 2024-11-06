package task

import (
	"context"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) (*models.Task, error)
	Get(ctx context.Context, id uint) (*models.Task, error)
	List(ctx context.Context, page int, limit int) ([]models.Task, int, int, int64, error)
	Update(ctx context.Context, task *models.Task) (*models.Task, error)
	Delete(ctx context.Context, id uint) (*models.Task, error)
}
