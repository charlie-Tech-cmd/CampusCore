package notification

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task defines the payload structure for our background processing queue
type Task struct {
	RecipientID string
	Type        string // e.g., "email_tuition_receipt", "email_course_registration"
	Payload     string
}

// Worker manages our concurrent background notification engine
type Worker struct {
	queue      chan Task
	wg         sync.WaitGroup
	shutdownChan chan struct{}
}

// NewWorker initializes a thread-safe notification buffer channel queue
func NewWorker(bufferSize int) *Worker {
	return &Worker{
		queue:        make(chan Task, bufferSize),
		shutdownChan: make(chan struct{}),
	}
}

// Start spawns a dedicated background goroutine worker pool to drain the task queue
func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		log.Println("📢 Notification engine background worker pool safely activated.")

		for {
			select {
			case task, open := <-w.queue:
				if !open {
					// The channel was explicitly closed during server teardown
					return
				}
				// Process the notification task defensively
				w.sendNotification(task)

			case <-w.shutdownChan:
				// Global termination catch signal caught
				return
			}
		}
	}()
}

// Enqueue puts a notification task onto our buffer without blocking the caller process
func (w *Worker) Enqueue(task Task) bool {
	select {
	case w.queue <- task:
		return true // Enqueued successfully
	default:
		log.Printf("⚠️ WARNING: Notification worker channel queue full! Dropping log task for %s", task.RecipientID)
		return false // Queue buffer saturated (Defensive drop to protect server RAM memory limits)
	}
}

// sendNotification simulates our external network integrations (SMTP Mail servers, SMS Gateways)
func (w *Worker) sendNotification(t Task) {
	// Simulate minor network transport overhead delay safely (e.g., calling SendGrid/Twilio API)
	time.Sleep(100 * time.Millisecond)
	
	log.Printf("📬 [ASYNC NOTIFICATION SENT] Target: %s | Type: %s | Details: %s", 
		t.RecipientID, t.Type, t.Payload,
	)
}

// Stop executes an orderly drainage and teardown cascade during application shutdown
func (w *Worker) Stop(ctx context.Context) {
	log.Println("⚠️ Initiating graceful shutdown sequence for the notification worker pool...")
	
	close(w.shutdownChan) // Signal loops to stop processing new inputs
	close(w.queue)        // Close the task stream channel

	// Create a channel to flag when waitgroup counters reach zero
	finishedChan := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(finishedChan)
	}()

	// Block processing until workers finish handling remaining tasks or context deadline expires
	select {
	case <-finishedChan:
		log.Println("✅ Notification background worker pool closed down cleanly with zero leaks.")
	case <-ctx.Done():
		log.Println("🚨 Coerced shutdown timeout expired: Notification worker terminated forcefully.")
	}
}