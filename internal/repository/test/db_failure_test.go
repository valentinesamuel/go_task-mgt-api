// internal/repository/test/db_failure_test.go
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

func TestDatabaseFailures(t *testing.T) {
	// Setup faulty DB
	db := testutils.SetupFaultyDB()
	repo := repository.NewTaskRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name      string
		operation func() error
	}{
		{
			name: "create failure",
			operation: func() error {
				task := &models.Task{
					Title:    "Failed Create",
					Priority: models.PriorityHigh,
					Status:   "todo",
				}
				_, err := repo.Create(ctx, task)
				return err
			},
		},
		{
			name: "get failure",
			operation: func() error {
				_, err := repo.Get(ctx, 1)
				return err
			},
		},
		{
			name: "list failure",
			operation: func() error {
				_, err := repo.List(ctx)
				return err
			},
		},
		{
			name: "update failure",
			operation: func() error {
				task := &models.Task{
					ID:       1,
					Title:    "Failed Update",
					Priority: models.PriorityHigh,
					Status:   "todo",
				}
				_, err := repo.Update(ctx, task)
				return err
			},
		},
		{
			name: "delete failure",
			operation: func() error {
				_, err := repo.Delete(ctx, 1)
				return err
			},
		},
		{
			name: "context timeout",
			operation: func() error {
				timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
				defer cancel()
				time.Sleep(time.Millisecond)
				_, err := repo.List(timeoutCtx)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			assert.Error(t, err, "Expected database operation to fail")
		})
	}
}
