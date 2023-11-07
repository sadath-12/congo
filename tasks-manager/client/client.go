package client

import (
	"encoding/json"
	"fmt"
	"log"
	task "sadath/tasks"
	"sync"

	"github.com/hibiken/asynq"
)

// List of queue names.
const (
	QueueNotifications = "notifications"
	QueueWebhooks      = "webhooks"
	email              = "email"
)

// client.go
func Client(wg *sync.WaitGroup) {

	defer wg.Done()
	client := asynq.NewClient(asynq.RedisClusterClientOpt{
		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})

	defer client.Close()

	t1, err := task.NewWelcomeEmailTask(42)
	if err != nil {
		log.Fatal(err)
	}

	t2, err := task.NewReminderEmailTask(42)
	if err != nil {
		log.Fatal(err)
	}

	// Process the task immediately.
	info, err := client.Enqueue(t1, asynq.Queue(email))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)

	info, err = client.Enqueue(t2, asynq.Queue(email))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)

	p1, err := json.Marshal(map[string]interface{}{"to": 123, "from": 456})
	task := asynq.NewTask("notifications:email", p1)
	res, err := client.Enqueue(task, asynq.Queue(QueueNotifications))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("successfully enqueued: %+v\n", res)

	// Create "webhooks:sync" task and enqueue it to "webhooks" queue.
	p2, err := json.Marshal(map[string]interface{}{"data": 123})
	task = asynq.NewTask("webhooks:sync", p2)
	res, err = client.Enqueue(task, asynq.Queue(QueueWebhooks))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("successfully enqueued: %+v\n", res)

}
