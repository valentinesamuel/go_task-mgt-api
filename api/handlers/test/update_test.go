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
	"testing"
)

// api/handlers/test/update_test.go
func TestUpdateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validTask := &models.Task{
		ID:       1,
		Title:    "Updated Task",
		Priority: models.PriorityHigh,
		Status:   "done",
	}

	tests := []struct {
		name             string
		taskID           string
		input            models.Task
		mockExpectations func(*mocks.MockTaskRepository)
		wantStatus       int
		wantError        string
	}{
		{
			name:   "valid update",
			taskID: "1",
			input: models.Task{
				Title:    "Updated Task",
				Priority: models.PriorityHigh,
				Status:   "done",
			},
			mockExpectations: func(m *mocks.MockTaskRepository) {
				m.On("Get", uint(1)).Return(validTask, nil)
				m.On("Update", mock.MatchedBy(func(t *models.Task) bool {
					return t.ID == 1
				})).Return(validTask, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "not found",
			taskID: "999",
			input: models.Task{
				Title:    "Not Found Task",
				Priority: models.PriorityHigh,
				Status:   "todo",
			},
			mockExpectations: func(m *mocks.MockTaskRepository) {
				m.On("Get", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantError:  "Task not found",
		},
		{
			name:             "invalid id",
			taskID:           "abc",
			input:            models.Task{},
			mockExpectations: func(m *mocks.MockTaskRepository) {},
			wantStatus:       http.StatusBadRequest,
			wantError:        "Invalid ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockTaskRepository)
			tt.mockExpectations(mockRepo)

			handler := handlers.NewTaskHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.taskID}}

			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdateTask(c)

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
