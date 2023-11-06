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
		fmt.Println("this is middleware")
		if !limiter.Allow() {
			return &RateLimitError{
				RetryIn: time.Duration(rand.Intn(10)) * time.Second,
			}
		}
		return nil
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
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 10,
			// If error is due to rate limit, don't count the error as a failure.
			IsFailure:      func(err error) bool { return !IsRateLimitError(err) },
			RetryDelayFunc: retryDelay,
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(middleFunc)
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
