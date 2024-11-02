package handlers_test

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/api/handlers"
	"github.com/valentinesamuel/go_task-mgt-api/internal/mocks"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
)

func TestGetTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		taskID     string
		mockReturn *models.Task
		mockError  error
		wantStatus int
		wantError  string
	}{
		{
			name:   "valid task",
			taskID: "1",
			mockReturn: &models.Task{
				ID:       1,
				Title:    "Test Task",
				Priority: "high",
				Status:   "todo",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id format",
			taskID:     "abc",
			wantStatus: http.StatusBadRequest,
			wantError:  "Invalid ID format",
		},
		{
			name:       "task not found",
			taskID:     "999",
			mockReturn: nil,
			mockError:  gorm.ErrRecordNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTaskRepository)

			if tt.taskID != "abc" {
				id, _ := strconv.ParseUint(tt.taskID, 10, 32)
				mockRepo.On("Get", uint(id)).Return(tt.mockReturn, tt.mockError)
			}

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.taskID}}

			handler.GetTask(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.wantError != "" {
				assert.Contains(t, response["error"], tt.wantError)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		mockReturn []models.Task
		mockError  error
		wantStatus int
		wantError  string
	}{
		{
			name: "successful list",
			mockReturn: []models.Task{
				{ID: 1, Title: "Task 1"},
				{ID: 2, Title: "Task 2"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "internal error",
			mockError:  errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "Failed to retrieve tasks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTaskRepository)
			mockRepo.On("List").Return(tt.mockReturn, tt.mockError)

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler.ListTasks(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantError == "" {
				var response []models.Task
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Len(t, response, len(tt.mockReturn))
			} else {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.wantError)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
