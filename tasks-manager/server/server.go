package server

import (
	"context"
	"encoding/json"
	"log"
	"sadath/tasks"
	"sync"

	"github.com/hibiken/asynq"
)

// workers.go
func Workers(wg *sync.WaitGroup) {
	defer wg.Done()
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("email:welcome", sendWelcomeEmail)
	mux.HandleFunc("email:reminder", sendReminderEmail)
	mux.HandleFunc("example_task", sendExample)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func sendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	c := make(chan error, 1)
	var p tasks.EmailTaskPayload
	go func() {
		c <- json.Unmarshal(t.Payload(), &p)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case _ = <-c:
		log.Printf(" [*] Send Welcome Email to User %d", p.UserID)
		return nil
	}
}

func sendExample(ctx context.Context, t *asynq.Task) error {
	c := make(chan error, 1)
	var p tasks.EmailTaskPayload
	go func() {
		c <- json.Unmarshal(t.Payload(), &p)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case _ = <-c:
		log.Printf(" [*] Send example Email to User %d", p.UserID)
		return nil
	}

}

func sendReminderEmail(ctx context.Context, t *asynq.Task) error {
	c := make(chan error, 1)
	var p tasks.EmailTaskPayload
	go func() {
		c <- json.Unmarshal(t.Payload(), &p)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case _ = <-c:
		log.Printf(" [*] Send Reminder Email to User %d", p.UserID)
		return nil
	}
}
