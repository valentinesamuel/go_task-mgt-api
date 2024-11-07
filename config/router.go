package config

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/valentinesamuel/go_task-mgt-api/internal/auth"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task"
	"github.com/valentinesamuel/go_task-mgt-api/pkg/middleware"
)

func SetupRouter(authHandler auth.AuthHandler, taskHandler task.TaskHandler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.CustomLogger(), gin.Recovery())

	r.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		logger.Debugln("Health check")
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", authHandler.RegisterUser)
		authRoute.POST("/login", authHandler.LoginUser)
		authRoute.POST("/logout", authHandler.LogoutUser)
	}

	taskRoute := r.Group("/tasks")
	{
		taskRoute.Use(middleware.AuthMiddleware())
		taskRoute.POST("/", taskHandler.CreateTask)
		taskRoute.GET("/:id", taskHandler.GetTask)
		taskRoute.GET("/", taskHandler.ListTasks)
		taskRoute.PUT("/:id", taskHandler.UpdateTask)
		taskRoute.DELETE("/:id", taskHandler.DeleteTask)
	}

	return r
}
