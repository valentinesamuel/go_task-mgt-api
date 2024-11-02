package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
	"testing"
)

func TestTaskRepository_Get(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)
	task := testutils.CreateTestTask(t, repo)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "existing task",
			id:      task.ID,
			wantErr: false,
		},
		{
			name:    "non-existent task",
			id:      999,
			wantErr: true,
		},
		{
			name:    "zero id",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, task.ID, got.ID)
				assert.Equal(t, task.Title, got.Title)
			}
		})
	}
}

func TestTaskRepository_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)

	tasks := []*models.Task{
		{Title: "Task 1", Priority: "high", Status: "todo"},
		{Title: "Task 2", Priority: "medium", Status: "in_progress"},
		{Title: "Task 3", Priority: "low", Status: "done"},
	}

	for _, task := range tasks {
		_, err := repo.Create(task)
		assert.NoError(t, err)
	}

	got, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, got, len(tasks))
	for i, task := range got {
		assert.NotZero(t, task.ID)
		assert.Equal(t, tasks[i].Title, task.Title)
	}
}
