package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
)

func TestTaskRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)

	tests := []struct {
		name    string
		task    *models.Task
		timeout time.Duration
		wantErr bool
	}{
		{
			name: "valid task",
			task: &models.Task{
				Title:    "Test Task",
				Priority: models.PriorityHigh,
				Status:   "todo",
			},
			timeout: time.Second * 5,
			wantErr: false,
		},
		{
			name:    "nil task",
			task:    nil,
			timeout: time.Second * 5,
			wantErr: true,
		},
		{
			name: "timeout context",
			task: &models.Task{
				Title:    "Timeout Task",
				Priority: models.PriorityHigh,
				Status:   "todo",
			},
			timeout: time.Nanosecond,
			wantErr: true,
		},
		{
			name: "invalid priority",
			task: &models.Task{
				Title:    "Invalid Priority",
				Priority: "invalid",
				Status:   "todo",
			},
			timeout: time.Second * 5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			created, err := repo.Create(ctx, tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, created)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, created)
				assert.NotZero(t, created.ID)
				assert.Equal(t, tt.task.Title, created.Title)
				assert.Equal(t, tt.task.Priority, created.Priority)
				assert.Equal(t, tt.task.Status, created.Status)
				assert.NotZero(t, created.CreatedAt)
				assert.NotZero(t, created.UpdatedAt)
			}
		})
	}
}
