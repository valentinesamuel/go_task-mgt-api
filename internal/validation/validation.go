package validation

import (
	"errors"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
)

func ValidateTask(task *models.Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	if err := ValidateTitle(task.Title); err != nil {
		return err
	}

	if err := ValidateDescription(task.Description); err != nil {
		return err
	}

	if err := ValidatePriority(task.Priority); err != nil {
		return err
	}

	if err := ValidateStatus(task.Status); err != nil {
		return err
	}

	return nil
}

func ValidateTitle(title string) error {
	if len(title) < 3 || len(title) > 100 {
		return errors.New("title must be between 3 and 100 characters")
	}
	return nil
}

func ValidateDescription(desc string) error {
	if len(desc) > 1000 {
		return errors.New("description must not exceed 1000 characters")
	}
	return nil
}

func ValidatePriority(p models.Priority) error {
	switch p {
	case models.PriorityLow, models.PriorityMedium, models.PriorityHigh:
		return nil
	default:
		return errors.New("invalid priority value")
	}
}

func ValidateStatus(s models.Status) error {
	switch s {
	case models.StatusTodo, models.StatusInProgress, models.StatusDone:
		return nil
	default:
		return errors.New("invalid status value")
	}
}
