package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
	"testing"
)

func TestTaskRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)

	tests := []struct {
		name    string
		task    *models.Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: &models.Task{
				Title:    "Test Task",
				Priority: "high",
				Status:   "todo",
			},
			wantErr: false,
		},
		{
			name:    "nil task",
			task:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.Create(tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, created)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, created)
				assert.NotZero(t, created.ID)
				assert.Equal(t, tt.task.Title, created.Title)
			}
		})
	}
}
