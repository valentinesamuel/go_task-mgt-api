package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/validation"
	"gorm.io/gorm"
)

type taskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &taskRepositoryImpl{db: db}
}

func (r *taskRepositoryImpl) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, errors.New("task is empty")
	}

	if err := validation.ValidateTask(task); err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (r *taskRepositoryImpl) Get(ctx context.Context, id uint) (*models.Task, error) {
	if id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var task models.Task
	err := r.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *taskRepositoryImpl) GetTasksByStatus(status models.Status) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("status = ?", status).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}

	return tasks, nil
}

func (r *taskRepositoryImpl) List(ctx context.Context, page int, limit int) ([]models.Task, int, int, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	var tasks []models.Task
	var total int64
	err := r.db.WithContext(ctx).Model(&models.Task{}).
		Count(&total).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, page, limit, total, nil
}

func (r *taskRepositoryImpl) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	if task == nil || task.ID == 0 {
		return nil, errors.New("invalid task or task id")
	}

	var exists models.Task
	err := r.db.WithContext(ctx).First(&exists, task.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task not found with id %d", task.ID)
		}
		return nil, fmt.Errorf("failed to verify task existence: %w", err)
	}

	err = r.db.WithContext(ctx).Save(task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

func (r *taskRepositoryImpl) Delete(ctx context.Context, id uint) (*models.Task, error) {
	if id == 0 {
		return nil, errors.New("invalid task id")
	}

	var task models.Task
	err := r.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task not found with id %d", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	err = r.db.WithContext(ctx).Delete(&task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}

	return &task, nil
}
