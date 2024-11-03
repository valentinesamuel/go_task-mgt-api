package task

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/valentinesamuel/go_task-mgt-api/internal/models"
	"github.com/valentinesamuel/go_task-mgt-api/pkg"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type taskHandlerImpl struct {
	repo TaskRepository
	mu   sync.RWMutex
}

func NewTaskHandler(repo TaskRepository) TaskHandler {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &taskHandlerImpl{
		repo: repo,
	}
}

var validate = validator.New()

// CreateTask creates a new task
// @Summary Create a new task
// @Description Creates a new task and stores it in the database
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task details"
// @Success 201 {object} pkg.SwaggerSuccessResponse "Task created successfully"
// @Failure 400 {object} pkg.SwaggerErrorResponse "Invalid request payload or validation error"
// @Failure 500 {object} pkg.SwaggerErrorResponse "Failed to create task"
// @Router /tasks [post]
func (h *taskHandlerImpl) CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest,
			pkg.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid request payload",
				err.Error()))
		return
	}

	if err := validate.Struct(&task); err != nil {
		validationErrors := pkg.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest,
			pkg.NewErrorResponse(
				http.StatusBadRequest,
				"Validation error",
				validationErrors))
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	created, err := h.repo.Create(ctx, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			pkg.NewErrorResponse(http.StatusInternalServerError,
				"Failed to create task",
				err.Error()))
		return
	}

	c.JSON(http.StatusCreated, pkg.NewSuccessResponse(
		http.StatusCreated,
		"Task created successfully",
		created,
	))
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Retrieve a task by its unique ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} pkg.SwaggerSuccessResponse
// @Failure 400 {object} pkg.SwaggerErrorResponse "Invalid ID format"
// @Failure 404 {object} pkg.SwaggerErrorResponse "Task not found"
// @Failure 500 {object} pkg.SwaggerErrorResponse "Failed to retrieve task"
// @Router /tasks/{id} [get]
func (h *taskHandlerImpl) GetTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, pkg.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err.Error()))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	task, err := h.repo.Get(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, pkg.NewErrorResponse(http.StatusNotFound, "Task not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, pkg.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve task", err.Error()))
		return
	}

	c.JSON(http.StatusOK, pkg.NewSuccessResponse(http.StatusOK, "Task retrieved successfully", task))
}

// ListTasks godoc
// @Summary List all tasks
// @Description Retrieve a list of all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} pkg.SwaggerSuccessResponse
// @Failure 500 {object} pkg.SwaggerErrorResponse "Failed to retrieve tasks"
// @Router /tasks [get]
func (h *taskHandlerImpl) ListTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	h.mu.RLock()
	defer h.mu.RUnlock()

	tasks, err := h.repo.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve tasks", err.Error()))
		return
	}

	c.JSON(http.StatusOK, pkg.NewSuccessResponse(http.StatusOK, "Tasks retrieved successfully", tasks))
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update an existing task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task data"
// @Success 200 {object} models.Task "Task updated successfully"
// @Failure 400 {object} pkg.SwaggerErrorResponse "Invalid ID format or validation error"
// @Failure 404 {object} pkg.SwaggerErrorResponse "Task not found"
// @Failure 500 {object} pkg.SwaggerErrorResponse "Failed to update task"
// @Router /tasks/{id} [put]
func (h *taskHandlerImpl) UpdateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, pkg.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err.Error()))
		return
	}

	h.mu.RLock()
	_, err = h.repo.Get(ctx, uint(id))
	h.mu.RUnlock()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, pkg.NewErrorResponse(http.StatusNotFound, "Task not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, pkg.NewErrorResponse(http.StatusInternalServerError, "Failed to retrieve task", err.Error()))
		return
	}

	var task models.Task
	if err := validate.Struct(&task); err != nil {
		validationErrors := pkg.FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, pkg.NewErrorResponse(http.StatusBadRequest, "Validation error", validationErrors))
		return
	}
	task.ID = uint(id)

	h.mu.Lock()
	defer h.mu.Unlock()

	updated, err := h.repo.Update(ctx, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pkg.NewErrorResponse(http.StatusInternalServerError, "Failed to update task", err.Error()))
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete an existing task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} pkg.SwaggerSuccessResponse "Task deleted successfully"
// @Failure 400 {object} pkg.SwaggerErrorResponse "Invalid ID format"
// @Failure 404 {object} pkg.SwaggerErrorResponse "Task not found"
// @Failure 500 {object} pkg.SwaggerErrorResponse "Failed to delete task"
// @Router /tasks/{id} [delete]
func (h *taskHandlerImpl) DeleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, pkg.NewErrorResponse(http.StatusBadRequest, "Invalid ID format", err.Error()))
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	deleted, err := h.repo.Delete(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, pkg.NewErrorResponse(http.StatusNotFound, "Task not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, pkg.NewErrorResponse(http.StatusInternalServerError, "Failed to delete task", err.Error()))
		return
	}

	c.JSON(http.StatusOK, pkg.NewSuccessResponse(http.StatusOK, "Task deleted successfully", deleted))
}
