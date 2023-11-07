package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sadath/tasks"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"golang.org/x/time/rate"
)

// Rate is 10 events/sec and permits burst of at most 30 events.
var limiter = rate.NewLimiter(10, 30)

func middleFunc(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		fmt.Println("you have my back ❤️ just move")
		if !limiter.Allow() {
			return &RateLimitError{
				RetryIn: time.Duration(rand.Intn(10)) * time.Second,
			}
		}
		return next.ProcessTask(ctx, t)
	})
}

type RateLimitError struct {
	RetryIn time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited (retry in  %v)", e.RetryIn)
}

func IsRateLimitError(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

func retryDelay(n int, err error, task *asynq.Task) time.Duration {
	var ratelimitErr *RateLimitError
	if errors.As(err, &ratelimitErr) {
		return ratelimitErr.RetryIn
	}
	return asynq.DefaultRetryDelayFunc(n, err, task)
}

// workers.go
func Workers(wg *sync.WaitGroup) {
	defer wg.Done()
	srv := asynq.NewServer(
		asynq.RedisClusterClientOpt{
			Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
		},
		asynq.Config{
			Concurrency: 10,
			// If error is due to rate limit, don't count the error as a failure.
			IsFailure:      func(err error) bool { return !IsRateLimitError(err) },
			RetryDelayFunc: retryDelay,
			Queues: map[string]int{
				"notifications": 1,
				"webhooks":      1,
				"email":         1,
				"example":       1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(middleFunc)
	mux.HandleFunc("email:welcome", sendWelcomeEmail)
	mux.HandleFunc("email:reminder", sendReminderEmail)
	mux.HandleFunc("example_task", sendExample)
	mux.HandleFunc("notifications:email", handleEmailTask)
	mux.HandleFunc("webhooks:sync", handleWebhookSyncTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func handleEmailTask(ctx context.Context, t *asynq.Task) error {
	c := make(chan error, 1)
	var p map[string]interface{}
	go func() {
		c <- json.Unmarshal(t.Payload(), &p)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		from, ok := p["from"].(float64)
		if !ok {
			fmt.Println("Value not a number")
			from = 0
		}

		to, ok := p["to"].(float64)
		if !ok {
			fmt.Println("Value not a number")
			to = 0
		}

		fmt.Printf("Send email from %d to %d\n", int(from), int(to))
		return nil
	}
}

func handleWebhookSyncTask(ctx context.Context, t *asynq.Task) error {
	c := make(chan error, 1)
	var p map[string]interface{}
	go func() {
		c <- json.Unmarshal(t.Payload(), &p)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case _ = <-c:
		data, ok := p["data"].(float64)
		if !ok {
			fmt.Println("Value not a number")
			data = 0
		}

		log.Printf(" [*] Send webhook to %d", int(data))
		return nil
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
