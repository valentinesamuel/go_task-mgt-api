package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valentinesamuel/go_task-mgt-api/api/handlers"
	"github.com/valentinesamuel/go_task-mgt-api/internal/mocks"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
)

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		input      interface{} // Changed to interface{} to handle invalid JSON
		mockReturn *models.Task
		mockError  error
		wantStatus int
		wantError  string
	}{
		{
			name: "valid task",
			input: models.Task{
				Title:    "Test Task",
				Priority: "high",
				Status:   "todo",
			},
			mockReturn: &models.Task{ID: 1, Title: "Test Task"},
			mockError:  nil,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid task",
			input:      models.Task{}, // Empty task
			wantStatus: http.StatusBadRequest,
			wantError:  "Key: 'Task.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		},
		{
			name: "internal server error",
			input: models.Task{
				Title:    "Test Task",
				Priority: "high",
				Status:   "todo",
			},
			mockReturn: nil,
			mockError:  errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTaskRepository)

			// Only set expectation if we expect repository to be called
			if tt.wantStatus != http.StatusBadRequest {
				mockRepo.On("Create", mock.AnythingOfType("*models.Task")).Return(tt.mockReturn, tt.mockError)
			}

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateTask(c)

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
