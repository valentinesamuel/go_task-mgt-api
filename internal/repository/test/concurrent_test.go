// internal/repository/test/concurrent_test.go
package repository_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentOperations(t *testing.T) {
	db := testutils.SetupTestDB(t)
	db.Exec("DELETE FROM tasks")
	repo := repository.NewTaskRepository(db)

	baseTasks := []models.Task{
		{Title: "Task 1", Priority: models.PriorityHigh, Status: models.StatusTodo},
		{Title: "Task 2", Priority: models.PriorityMedium, Status: models.StatusInProgress},
	}

	ctx := context.Background()
	for _, task := range baseTasks {
		_, err := repo.Create(ctx, &task)
		assert.NoError(t, err)
	}

	var successfulWrites int32
	operationCount := 10
	var wg sync.WaitGroup

	for i := 0; i < operationCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			task := &models.Task{
				Title:    fmt.Sprintf("Concurrent Task %d", i),
				Priority: models.PriorityHigh,
				Status:   models.StatusTodo,
			}
			_, err := repo.Create(ctx, task)
			if err == nil {
				atomic.AddInt32(&successfulWrites, 1)
			}
		}(i)
	}

	wg.Wait()

	tasks, err := repo.List(ctx)
	assert.NoError(t, err)

	expectedCount := int(atomic.LoadInt32(&successfulWrites)) + len(baseTasks)
	assert.Len(t, tasks, expectedCount, "Expected exactly %d tasks but got %d", expectedCount, len(tasks))

	for _, task := range tasks {
		assert.NotZero(t, task.ID)
		assert.NotEmpty(t, task.Title)
		assert.NotEmpty(t, task.Priority)
		assert.NotEmpty(t, task.Status)
	}
}
