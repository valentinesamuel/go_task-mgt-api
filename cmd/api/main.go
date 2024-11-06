package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/valentinesamuel/go_task-mgt-api/config"
	_ "github.com/valentinesamuel/go_task-mgt-api/docs"
	"github.com/valentinesamuel/go_task-mgt-api/internal/auth"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task"
	"github.com/valentinesamuel/go_task-mgt-api/internal/user"
	"github.com/valentinesamuel/go_task-mgt-api/pkg/middleware"
	"log"
)

// @title			Task Management API
// @version		1.0
// @description	API for managing tasks and users in the system
// @host			localhost:8080
// @BasePath		/
func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	authRepo := user.NewUserRepository(db)
	authHandler := auth.NewAuthHandler(authRepo)

	taskRepo := task.NewTaskRepository(db)
	taskHandler := task.NewTaskHandler(taskRepo)

	r := gin.Default()

	r.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
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

	log.Printf("Starting server on port 8080. Visit the docs at http://localhost:8080/docs/index.html")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
