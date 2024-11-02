package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
	"testing"
)

func TestTaskRepository_Delete(t *testing.T) {
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
			deleted, err := repo.Delete(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, deleted)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, deleted)

				// Verify deletion
				_, err := repo.Get(tt.id)
				assert.Error(t, err)
			}
		})
	}
}
