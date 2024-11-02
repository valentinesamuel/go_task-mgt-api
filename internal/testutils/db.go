package testutils

import (
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func SetupTestDB(t *testing.T) *gorm.DB {
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

func CreateTestTask(t *testing.T, repo repository.TaskRepository) *models.Task {
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
