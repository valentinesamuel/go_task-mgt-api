package repository

import (
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
)

type TaskRepository interface {
	Create(task *models.Task) (*models.Task, error)
	Update(task *models.Task) (*models.Task, error)
	Delete(id uint) (*models.Task, error)
	Get(id uint) (*models.Task, error)
	List() ([]models.Task, error)
}
