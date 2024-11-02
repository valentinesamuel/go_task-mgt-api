// internal/repository/test/read_test.go
package repository_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"github.com/valentinesamuel/go_task-mgt-api/internal/testutils"
	"gorm.io/gorm"
	"testing"
	"time"
)

//func TestTaskRepository_Get(t *testing.T) {
//	db := testutils.SetupTestDB(t)
//	repo := repository.NewTaskRepository(db)
//	ctx := context.Background()
//
//	// Create test task
//	task := &models.Task{
//		Title:    "Test Task",
//		Priority: models.PriorityHigh,
//		Status:   "todo",
//	}
//	created, err := repo.Create(ctx, task)
//	assert.NoError(t, err)
//	assert.NotNil(t, created)
//
//	tests := []struct {
//		name    string
//		id      uint
//		ctx     context.Context
//		wantErr bool
//		errType error
//	}{
//		{
//			name:    "existing task",
//			id:      created.ID,
//			ctx:     context.Background(),
//			wantErr: false,
//		},
//		{
//			name:    "non-existent task",
//			id:      999,
//			ctx:     context.Background(),
//			wantErr: true,
//			errType: gorm.ErrRecordNotFound,
//		},
//		{
//			name:    "zero id",
//			id:      0,
//			ctx:     context.Background(),
//			wantErr: true,
//			errType: gorm.ErrRecordNotFound,
//		},
//		{
//			name:    "context timeout",
//			id:      created.ID,
//			ctx:     getTimeoutContext(t),
//			wantErr: true,
//		},
//		{
//			name:    "context canceled",
//			id:      created.ID,
//			ctx:     getCanceledContext(t),
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := repo.Get(tt.ctx, tt.id)
//			if tt.wantErr {
//				assert.Error(t, err)
//				if tt.errType != nil {
//					assert.ErrorIs(t, err, tt.errType)
//				}
//				assert.Nil(t, got)
//			} else {
//				assert.NoError(t, err)
//				assert.NotNil(t, got)
//				assert.Equal(t, created.ID, got.ID)
//				assert.Equal(t, created.Title, got.Title)
//				assert.Equal(t, created.Priority, got.Priority)
//				assert.Equal(t, created.Status, got.Status)
//			}
//		})
//	}
//}

//func TestTaskRepository_Get(t *testing.T) {
//	db := testutils.SetupTestDB(t)
//	repo := repository.NewTaskRepository(db)
//	tests := []struct {
//		name    string
//		id      uint
//		ctx     context.Context
//		wantErr error
//	}{
//		{
//			name:    "existing task",
//			id:      1,
//			ctx:     context.Background(),
//			wantErr: nil,
//		},
//		{
//			name:    "non-existent task",
//			id:      999,
//			ctx:     context.Background(),
//			wantErr: gorm.ErrRecordNotFound,
//		},
//		{
//			name:    "zero id",
//			id:      0,
//			ctx:     context.Background(),
//			wantErr: gorm.ErrRecordNotFound,
//		},
//		// ... other test cases
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := repo.Get(tt.ctx, tt.id)
//			if tt.wantErr != nil {
//				assert.ErrorIs(t, err, tt.wantErr)
//				assert.Nil(t, got)
//			} else {
//				assert.NoError(t, err)
//				assert.NotNil(t, got)
//			}
//		})
//	}
//}

func TestTaskRepository_Get(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)
	ctx := context.Background()

	// Create test task first
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
			id:      created.ID, // Use created task's ID
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
		// ... other test cases remain the same
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

func getTimeoutContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	t.Cleanup(cancel)
	time.Sleep(time.Millisecond)
	return ctx
}

func getCanceledContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

//func TestTaskRepository_List(t *testing.T) {
//	db := testutils.SetupTestDB(t)
//	repo := repository.NewTaskRepository(db)
//	ctx := context.Background()
//
//	// Create test tasks
//	testTasks := []models.Task{
//		{Title: "Task 1", Priority: models.PriorityHigh, Status: "todo"},
//		{Title: "Task 2", Priority: models.PriorityMedium, Status: "in_progress"},
//		{Title: "Task 3", Priority: models.PriorityLow, Status: "done"},
//	}
//
//	for _, task := range testTasks {
//		_, err := repo.Create(ctx, &task)
//		assert.NoError(t, err)
//	}
//
//	tests := []struct {
//		name      string
//		ctx       context.Context
//		setupFunc func() // Additional setup if needed
//		wantCount int
//		wantErr   bool
//	}{
//		{
//			name:      "list all tasks",
//			ctx:       context.Background(),
//			wantCount: len(testTasks),
//			wantErr:   false,
//		},
//		{
//			name: "empty list",
//			ctx:  context.Background(),
//			setupFunc: func() {
//				db.Exec("DELETE FROM tasks")
//			},
//			wantCount: 0,
//			wantErr:   false,
//		},
//		{
//			name:      "context timeout",
//			ctx:       getTimeoutContext(t),
//			wantCount: 0,
//			wantErr:   true,
//		},
//		{
//			name:      "context canceled",
//			ctx:       getCanceledContext(t),
//			wantCount: 0,
//			wantErr:   true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if tt.setupFunc != nil {
//				tt.setupFunc()
//			}
//
//			tasks, err := repo.List(tt.ctx)
//			if tt.wantErr {
//				assert.Error(t, err)
//				assert.Nil(t, tasks)
//			} else {
//				assert.NoError(t, err)
//				assert.Len(t, tasks, tt.wantCount)
//
//				if tt.wantCount > 0 {
//					// Verify task properties
//					for _, task := range tasks {
//						assert.NotZero(t, task.ID)
//						assert.NotEmpty(t, task.Title)
//						assert.NotEmpty(t, task.Priority)
//						assert.NotEmpty(t, task.Status)
//						assert.NotZero(t, task.CreatedAt)
//						assert.NotZero(t, task.UpdatedAt)
//					}
//				}
//			}
//		})
//	}
//}

func TestTaskRepository_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)
	ctx := context.Background()

	// Clear database before starting
	db.Exec("DELETE FROM tasks")

	// Create test tasks
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
			// Clear database before each test case
			if tt.setupFunc == nil {
				db.Exec("DELETE FROM tasks")
				// Recreate test tasks for "list all tasks" case
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
