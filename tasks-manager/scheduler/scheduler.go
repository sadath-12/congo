package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/hibiken/asynq"
)

func handleEnqueueError(_ *asynq.Task, _ []asynq.Option, _ error) {
	log.Println("Handling the errors just for you ❤️")
}

func Scheduler(wg *sync.WaitGroup) {
	defer wg.Done()

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		&asynq.SchedulerOpts{
			Location:            loc,
			EnqueueErrorHandler: handleEnqueueError,
		},
	)

	task := asynq.NewTask("example_task", nil)

	// You can use "@every <duration>" to specify the interval.
	entryID, err := scheduler.Register("@every 1s", task)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
