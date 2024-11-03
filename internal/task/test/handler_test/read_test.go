package handler_test_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task/test/repository_test"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/api/handlers"
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
				Priority: models.PriorityHigh,
				Status:   models.StatusTodo,
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
			mockError:  gorm.ErrRecordNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository_test.MockTaskRepository)

			if tt.taskID != "abc" {
				id, _ := strconv.ParseUint(tt.taskID, 10, 32)
				mockRepo.On("Get", mock.Anything, uint(id)).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request with proper context
			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.taskID, nil)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: tt.taskID}}

			handler.GetTask(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.wantError != "" {
				assert.Contains(t, response["error"], tt.wantError)
			} else {
				assert.Equal(t, float64(tt.mockReturn.ID), response["id"])
				assert.Equal(t, tt.mockReturn.Title, response["title"])
				assert.Equal(t, string(tt.mockReturn.Priority), response["priority"])
				assert.Equal(t, string(tt.mockReturn.Status), response["status"])
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
				{
					ID:       1,
					Title:    "Task 1",
					Priority: models.PriorityHigh,
					Status:   models.StatusTodo,
				},
				{
					ID:       2,
					Title:    "Task 2",
					Priority: models.PriorityMedium,
					Status:   models.StatusInProgress,
				},
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
			mockRepo := new(repository_test.MockTaskRepository)
			mockRepo.On("List", mock.Anything).Return(tt.mockReturn, tt.mockError).Once()

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request with proper context
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			c.Request = req

			handler.ListTasks(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					return
				}
				assert.Contains(t, response["error"], tt.wantError)
			} else {
				var response []models.Task
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					return
				}
				assert.Equal(t, len(tt.mockReturn), len(response))
				for i, task := range response {
					assert.Equal(t, tt.mockReturn[i].ID, task.ID)
					assert.Equal(t, tt.mockReturn[i].Title, task.Title)
					assert.Equal(t, tt.mockReturn[i].Priority, task.Priority)
					assert.Equal(t, tt.mockReturn[i].Status, task.Status)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
