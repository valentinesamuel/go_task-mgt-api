package repository

import (
	"errors"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"gorm.io/gorm"
)

type taskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepositoryImpl{
		db: db,
	}
}

func (r *taskRepositoryImpl) Create(task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, errors.New("task is empty")
	}

	if err := r.db.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (r *taskRepositoryImpl) Get(id uint) (*models.Task, error) {
	if id == 0 {
		return nil, errors.New("task not found")
	}

	task := models.Task{}
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskRepositoryImpl) List() ([]models.Task, error) {
	var tasks []models.Task

	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepositoryImpl) Update(task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, errors.New("task is invalid")
	}

	var exists models.Task
	if err := r.db.First(&exists, task.ID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Save(task).Error; err != nil {
		return nil, err
	}

	var updated models.Task
	if err := r.db.First(&updated, task.ID).Error; err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *taskRepositoryImpl) Delete(id uint) (*models.Task, error) {
	if id == 0 {
		return nil, errors.New("task not found")
	}

	var task models.Task
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Delete(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}
