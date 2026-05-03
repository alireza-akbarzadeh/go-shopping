package jobs

import (
	"log"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/tasks"
	"github.com/robfig/cron/v3"
)

type CronJobs struct {
	scheduler *cron.Cron
	taskPool  *tasks.WorkerPool
	svc       *services.Services
}

// Recoverer is a job wrapper that recovers from panics to prevent the scheduler from crashing.
func Recoverer(next cron.Job) cron.Job {
	return cron.FuncJob(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Cron job panic recovered: %v", r)
			}
		}()
		next.Run()
	})
}

func NewCronJobs(taskPool *tasks.WorkerPool, svc *services.Services) *CronJobs {
	scheduler := cron.New(
		cron.WithLocation(time.UTC),
		cron.WithChain(Recoverer),
	)
	return &CronJobs{
		scheduler: scheduler,
		taskPool:  taskPool,
		svc:       svc,
	}
}

// registerJobs centralizes all your scheduled tasks.
func (c *CronJobs) registerJobs() {
	c.registerCartJobs()
	c.registerProductJobs()
	c.registerOrderJobs()
}

// addJob is a helper to add a job to the scheduler with error logging.
func (c *CronJobs) addJob(schedule, name string, cmd func()) {
	_, err := c.scheduler.AddFunc(schedule, cmd)
	if err != nil {
		log.Printf("Failed to schedule job '%c' with schedule '%c': %v", name, schedule, err)
	}
}

// Start --------------------------------
// Lifecycle
// --------------------------------
func (c *CronJobs) Start() {
	c.registerJobs()
	c.scheduler.Start()
	log.Println("CronJobs started.")
}

func (c *CronJobs) Stop() {
	ctx := c.scheduler.Stop()
	<-ctx.Done()
	log.Println("CronJobs stopped.")
}
