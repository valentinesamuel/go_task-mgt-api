package validation

import (
	"github.com/stretchr/testify/assert"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"strings"
	"testing"
)

func TestValidateTask(t *testing.T) {
	tests := []struct {
		name    string
		task    *models.Task
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid task",
			task: &models.Task{
				Title:       "Valid Task",
				Description: "Valid description",
				Priority:    models.PriorityHigh,
				Status:      "todo",
			},
			wantErr: false,
		},
		{
			name:    "nil task",
			task:    nil,
			wantErr: true,
			errMsg:  "task cannot be nil",
		},
		{
			name: "title too short",
			task: &models.Task{
				Title:       "ab",
				Description: "Valid description",
				Priority:    models.PriorityHigh,
				Status:      "todo",
			},
			wantErr: true,
			errMsg:  "title must be between 3 and 100 characters",
		},
		{
			name: "title too long",
			task: &models.Task{
				Title:       strings.Repeat("a", 101),
				Description: "Valid description",
				Priority:    models.PriorityHigh,
				Status:      "todo",
			},
			wantErr: true,
			errMsg:  "title must be between 3 and 100 characters",
		},
		{
			name: "description too long",
			task: &models.Task{
				Title:       "Valid Title",
				Description: strings.Repeat("a", 1001),
				Priority:    models.PriorityHigh,
				Status:      "todo",
			},
			wantErr: true,
			errMsg:  "description must not exceed 1000 characters",
		},
		{
			name: "invalid priority",
			task: &models.Task{
				Title:       "Valid Title",
				Description: "Valid description",
				Priority:    "invalid",
				Status:      "todo",
			},
			wantErr: true,
			errMsg:  "invalid priority value",
		},
		{
			name: "invalid status",
			task: &models.Task{
				Title:       "Valid Title",
				Description: "Valid description",
				Priority:    models.PriorityHigh,
				Status:      "invalid",
			},
			wantErr: true,
			errMsg:  "invalid status value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTask(tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
