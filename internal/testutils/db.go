package testutils

import (
	"context"
	"database/sql"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func CreateTestTask(t *testing.T, repo task.TaskRepository) *models.Task {
	ctx := context.Background() // Add context
	testTask := &models.Task{
		Title:    "Test Task",
		Priority: "high",
		Status:   "todo",
	}
	created, err := repo.Create(ctx, testTask) // Pass context to Create
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}
	return created
}

func SetupFaultyDB() *gorm.DB {
	sqlDB, _ := sql.Open("sqlite3", ":memory:")
	_ = sqlDB.Close() // Close immediately to simulate failure
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	return db
}
