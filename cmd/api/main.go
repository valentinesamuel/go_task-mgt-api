package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/valentinesamuel/go_task-mgt-api/config"
	_ "github.com/valentinesamuel/go_task-mgt-api/docs"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task"
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

	taskRepo := task.NewTaskRepository(db)
	taskHandler := task.NewTaskHandler(taskRepo)

	r := gin.Default()

	r.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	tasks := r.Group("/tasks")
	{
		tasks.POST("/", taskHandler.CreateTask)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.GET("/", taskHandler.ListTasks)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)

	}

	log.Printf("Starting server on port 8080. Visit the docs at http://localhost:8080/docs/index.html")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
