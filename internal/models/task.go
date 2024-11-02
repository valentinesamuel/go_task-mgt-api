package models

import (
	"time"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Title       string    `json:"title" binding:"required,min=3,max=100"`
	Description string    `json:"description" binding:"max=100"`
	Priority    Priority  `json:"priority" binding:"required,oneof=low medium high"`
	Status      string    `json:"status" binding:"required,oneof=todo in_progress done"`
	DueDate     time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
