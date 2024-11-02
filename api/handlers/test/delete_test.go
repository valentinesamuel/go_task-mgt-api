package handlers_test

import (
	"encoding/json"
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

// api/handlers/test/delete_test.go
func TestDeleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		taskID        string
		mockGetReturn *models.Task
		mockGetError  error
		mockDelReturn *models.Task
		mockDelError  error
		wantStatus    int
		wantError     string
	}{
		{
			name:          "successful delete",
			taskID:        "1",
			mockGetReturn: &models.Task{ID: 1},
			mockGetError:  nil,
			mockDelReturn: &models.Task{ID: 1},
			mockDelError:  nil,
			wantStatus:    http.StatusOK,
		},
		{
			name:       "invalid id",
			taskID:     "abc",
			wantStatus: http.StatusBadRequest,
			wantError:  "Invalid ID format",
		},
		{
			name:          "not found",
			taskID:        "999",
			mockGetReturn: nil,
			mockGetError:  gorm.ErrRecordNotFound,
			wantStatus:    http.StatusNotFound,
			wantError:     "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTaskRepository)

			if tt.taskID != "abc" {
				id, _ := strconv.ParseUint(tt.taskID, 10, 32)
				mockRepo.On("Get", uint(id)).Return(tt.mockGetReturn, tt.mockGetError)
				if tt.mockGetError == nil {
					mockRepo.On("Delete", uint(id)).Return(tt.mockDelReturn, tt.mockDelError)
				}
			}

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.taskID}}

			handler.DeleteTask(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			if tt.wantError != "" {
				assert.Contains(t, response["error"], tt.wantError)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
