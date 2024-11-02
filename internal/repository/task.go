package repository

import (
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) Create(task *models.Task) error {
	result := r.db.Create(task)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TaskRepository) Get(id uint) (*models.Task, error) {
	var task models.Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func (r *TaskRepository) List() ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	result := r.db.Save(&task)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TaskRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
