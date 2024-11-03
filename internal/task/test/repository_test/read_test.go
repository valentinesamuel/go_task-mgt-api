package repository_test_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
	"gorm.io/gorm"
	"testing"
)

func TestTaskRepository_Get(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)
	ctx := context.Background()

	testTask := &models.Task{
		Title:    "Test Task",
		Priority: models.PriorityHigh,
		Status:   models.StatusTodo,
	}
	created, err := repo.Create(ctx, testTask)
	assert.NoError(t, err)
	assert.NotNil(t, created)

	tests := []struct {
		name    string
		id      uint
		ctx     context.Context
		wantErr error
	}{
		{
			name:    "existing task",
			id:      created.ID,
			ctx:     context.Background(),
			wantErr: nil,
		},
		{
			name:    "non-existent task",
			id:      999,
			ctx:     context.Background(),
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "zero id",
			id:      0,
			ctx:     context.Background(),
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(tt.ctx, tt.id)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.id, got.ID)
			}
		})
	}
}

func TestTaskRepository_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)
	ctx := context.Background()

	db.Exec("DELETE FROM tasks")

	testTasks := []models.Task{
		{Title: "Task 1", Priority: models.PriorityHigh, Status: models.StatusTodo},
		{Title: "Task 2", Priority: models.PriorityMedium, Status: models.StatusInProgress},
		{Title: "Task 3", Priority: models.PriorityLow, Status: models.StatusDone},
	}

	for _, task := range testTasks {
		_, err := repo.Create(ctx, &task)
		assert.NoError(t, err)
	}

	var tests []struct {
		name      string
		ctx       context.Context
		setupFunc func()
		wantCount int
		wantErr   bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc == nil {
				db.Exec("DELETE FROM tasks")
				if tt.name == "list all tasks" {
					for _, task := range testTasks {
						_, err := repo.Create(ctx, &task)
						assert.NoError(t, err)
					}
				}
			} else {
				tt.setupFunc()
			}

		})
	}
}
