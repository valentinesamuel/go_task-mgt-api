package workers

import (
	"github.com/robfig/cron/v3"
	"github.com/valentinesamuel/go_task-mgt-api/internal/task"
	"log"
	"time"
)

type CronService struct {
	cron     *cron.Cron
	taskRepo task.TaskRepository
}

func NewCronService(taskRepo task.TaskRepository) *CronService {
	return &CronService{
		cron:     cron.New(cron.WithSeconds()),
		taskRepo: taskRepo,
	}
}

func (cs *CronService) StartCron() {
	// Check in-progress tasks every minute
	_, err := cs.cron.AddFunc("* * * * * *", func() {
		tasks, err := cs.taskRepo.GetTasksByStatus("in_progress")
		if err != nil {
			log.Printf("Error fetching in-progress tasks: %v", err)
			return
		}

		log.Printf("In-progress tasks at %v:", time.Now().Format(time.RFC3339))
		for _, t := range tasks {
			log.Printf("Task ID: %d, Title: %s, Status: %s",
				t.ID, t.Title, t.Status)
		}
	})

	if err != nil {
		log.Fatal("Error setting up cron job:", err)
	}

	cs.cron.Start()
	log.Println("Cron service started successfully")
}

func (cs *CronService) StopCron() {
	cs.cron.Stop()
}
