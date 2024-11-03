package handler_test_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valentinesamuel/go_task-mgt-api/api/handlers"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task/test/repository_test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		input      interface{}
		setupCtx   func(*gin.Context)
		mockReturn *models.Task
		mockError  error
		wantStatus int
		wantError  string
	}{
		{
			name: "valid task",
			input: models.Task{
				Title:    "Test Task",
				Priority: models.PriorityHigh,
				Status:   models.StatusTodo, // Use Status type
			},
			mockReturn: &models.Task{
				ID:       1,
				Title:    "Test Task",
				Priority: models.PriorityHigh,
				Status:   models.StatusTodo, // Use Status type
			},
			wantStatus: http.StatusCreated,
		},
		// ... other test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository_test.MockTaskRepository)

			// Fix mock expectation to handle context
			if tt.wantStatus != http.StatusBadRequest {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Task")).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			if tt.setupCtx != nil {
				tt.setupCtx(c)
			}

			handler.CreateTask(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				return
			}

			if tt.wantError != "" {
				assert.Contains(t, response["error"], tt.wantError)
			} else if tt.mockReturn != nil {
				assert.Equal(t, float64(tt.mockReturn.ID), response["id"])
				assert.Equal(t, tt.mockReturn.Title, response["title"])
				assert.Equal(t, string(tt.mockReturn.Priority), response["priority"])
				assert.Equal(t, string(tt.mockReturn.Status), response["status"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
