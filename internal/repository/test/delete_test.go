package repository_test

import (
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

func TestDeleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		taskID     string
		setupCtx   func(*gin.Context)
		mockReturn *models.Task
		mockError  error
		wantStatus int
		wantError  string
	}{
		{
			name:   "successful delete",
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
			name:       "invalid id",
			taskID:     "abc",
			wantStatus: http.StatusBadRequest,
			wantError:  "Invalid ID format",
		},
		{
			name:       "not found",
			taskID:     "999",
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
				mockRepo.On("Delete", mock.Anything, uint(id)).
					Return(tt.mockReturn, tt.mockError).
					Once()
			}

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request with proper context
			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.taskID, nil)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: tt.taskID}}

			if tt.setupCtx != nil {
				tt.setupCtx(c)
			}

			handler.DeleteTask(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				return
			}

			if tt.wantError != "" {
				assert.Contains(t, response["error"], tt.wantError)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
