package repository

import (
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func createTestTask(t *testing.T, repo TaskRepository) *models.Task {
	task := &models.Task{
		Title:    "Test Task",
		Priority: "high",
		Status:   "todo",
	}
	created, err := repo.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	return created
}

func TestTaskRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

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
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && created.ID == 0 {
				t.Error("Expected ID to be set")
			}
		})
	}
}

func TestTaskRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)
	task := createTestTask(t, repo)

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
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("Get() got = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestTaskRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)
	task := createTestTask(t, repo)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if updated.Title != tt.task.Title {
					t.Errorf("Update() got = %v, want %v", updated.Title, tt.task.Title)
				}
			}
		})
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)
	task := createTestTask(t, repo)

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
			name:    "zero id",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify deletion
				_, err := repo.Get(tt.id)
				if err == nil {
					t.Error("Delete() task still exists")
				}
			}
		})
	}
}

func TestTaskRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	// Create multiple tasks
	for i := 0; i < 3; i++ {
		createTestTask(t, repo)
	}

	tasks, err := repo.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
		return
	}

	if len(tasks) != 3 {
		t.Errorf("List() got = %v items, want %v items", len(tasks), 3)
	}
}
