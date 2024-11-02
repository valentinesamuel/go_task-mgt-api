package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
	"testing"
)

func TestTaskRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)
	task := testutils.CreateTestTask(t, repo)

	tests := []struct {
		name    string
		task    *models.Task
		wantErr bool
	}{
		{
			name: "valid update",
			task: &models.Task{
				ID:       task.ID,
				Title:    "Updated Task",
				Priority: "low",
				Status:   "done",
			},
			wantErr: false,
		},
		{
			name:    "nil task",
			task:    nil,
			wantErr: true,
		},
		{
			name: "non-existent task",
			task: &models.Task{
				ID:    999,
				Title: "Non-existent",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, updated)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updated)
				assert.Equal(t, tt.task.Title, updated.Title)
				assert.Equal(t, tt.task.Priority, updated.Priority)
				assert.Equal(t, tt.task.Status, updated.Status)
			}
		})
	}
}
