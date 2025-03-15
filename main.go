package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/valentinesamuel/go_task-mgt-api/config"
	_ "github.com/valentinesamuel/go_task-mgt-api/docs"
	"github.com/valentinesamuel/go_task-mgt-api/internal/auth"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task"
	"github.com/valentinesamuel/go_task-mgt-api/internal/user"
	"os"
)

// @title			Task Management API
// @version		1.0
// @description	API for managing tasks and users in the system
// @host			localhost:8080
// @BasePath		/
func main() {
	file := config.SetupLogger()
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Fatal(err, "Failed to close the log file")
		}
	}(file)

	db := config.SetupDatabase()

	authRepo := user.NewUserRepository(db)
	authHandler := auth.NewAuthHandler(authRepo)

	taskRepo := task.NewTaskRepository(db)
	taskHandler := task.NewTaskHandler(taskRepo)

	//cronService := workers.NewCronService(taskRepo)
	//cronService.StartCron()

	r := config.SetupRouter(authHandler, taskHandler)

	logger.Info("Starting server on port 8080. Visit the docs at http://localhost:8080/docs/index.html")
	if err := r.Run(":8080"); err != nil {
		logger.Fatal(err)
	}
}
