package main

import (
	"github.com/gin-gonic/gin"
	"github.com/valentinesamuel/go_task-mgt-api/config"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task"
	"log"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	taskRepo := task.NewTaskRepository(db)
	taskHandler := task.NewTaskHandler(taskRepo)

	r := gin.Default()

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

	log.Printf("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
