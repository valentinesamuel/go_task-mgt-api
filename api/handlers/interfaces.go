package handlers

import (
	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	CreateTask(c *gin.Context)
	GetTask(c *gin.Context)
	ListTasks(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}
