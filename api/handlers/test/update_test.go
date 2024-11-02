package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valentinesamuel/go_task-mgt-api/api/handlers"
	"github.com/valentinesamuel/go_task-mgt-api/internal/mocks"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// api/handlers/test/update_test.go
// api/handlers/test/update_test.go
func TestUpdateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		taskID           string
		input            models.Task
		mockGetReturn    *models.Task
		mockGetError     error
		mockUpdateReturn *models.Task
		mockUpdateError  error
		wantStatus       int
		wantError        string
	}{
		{
			name:   "valid update",
			taskID: "1",
			input: models.Task{
				Title:    "Updated Task",
				Priority: models.PriorityHigh,
				Status:   models.StatusTodo,
			},
			mockGetReturn: &models.Task{ID: 1}, // For existence check
			mockUpdateReturn: &models.Task{
				ID:       1,
				Title:    "Updated Task",
				Priority: models.PriorityHigh,
				Status:   models.StatusTodo,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			taskID:     "abc",
			wantStatus: http.StatusBadRequest,
			wantError:  "Invalid ID format",
		},
		{
			name:   "not found",
			taskID: "999",
			input: models.Task{
				Title:    "Not Found Task",
				Priority: models.PriorityHigh,
				Status:   models.StatusTodo,
			},
			mockGetError: gorm.ErrRecordNotFound,
			wantStatus:   http.StatusNotFound,
			wantError:    "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTaskRepository)

			if tt.taskID != "abc" {
				id, _ := strconv.ParseUint(tt.taskID, 10, 32)
				// Setup Get mock
				mockRepo.On("Get", mock.Anything, uint(id)).
					Return(tt.mockGetReturn, tt.mockGetError).
					Once()

				if tt.mockGetError == nil {
					// Setup Update mock
					mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Task")).
						Return(tt.mockUpdateReturn, tt.mockUpdateError).
						Once()
				}
			}

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request with proper context
			req := httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, nil)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: tt.taskID}}

			// Add request body
			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdateTask(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.wantError != "" {
				assert.Contains(t, response["error"], tt.wantError)
			} else {
				assert.Equal(t, float64(tt.mockUpdateReturn.ID), response["id"])
				assert.Equal(t, tt.mockUpdateReturn.Title, response["title"])
				assert.Equal(t, string(tt.mockUpdateReturn.Priority), response["priority"])
				assert.Equal(t, string(tt.mockUpdateReturn.Status), response["status"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
