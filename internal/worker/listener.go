package worker

import (
	"container/heap"
	"context"
	"log"

	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartPgListener(ctx context.Context, channel string, pool *pgxpool.Pool, c *bootstrap.Container) {
	conn, _ := pool.Acquire(ctx)
	defer conn.Release()
	_, err := conn.Exec(ctx, `LISTEN `+channel)
	if err != nil {
		log.Fatal("[LISTENER] LISTEN failed:", err)
	}

	log.Println("[LISTENER] Listening on channel:", channel)
	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			log.Println("[LISTENER] Error while waiting:", err)
			continue
		}

		jID := notification.Payload
		log.Println("[LISTENER] Received notification for job ID:", jID)
		jobID, err := uuid.Parse(jID)
		if err != nil {
			log.Printf("[LISTENER] Invalid job ID format: %s\n", jobID)
			continue
		}

		jobPayload, err := c.JobHandler.GetJobPayload(ctx, jobID)
		if err != nil {
			log.Printf("[LISTENER] Failed to fetch job payload for ID %s: %v\n", jobID, err)
			continue
		}

		priorityValue := map[string]int{
			"high":   1,
			"medium": 2,
			"low":    3,
		}[jobPayload.Priority]

		job := &JobItem{
			ID:       jobID,
			RunAt:    jobPayload.RunAt,
			Priority: priorityValue,
			Attempts: jobPayload.Attempts,
			Payload:  jobPayload.Payload,
			JobType:  jobPayload.JobType,
			Status:   jobPayload.Status,
		}

		log.Printf("[LISTENER] Enqueuing job ID %s with priority %d and run_at %s\n", jobID, priorityValue, jobPayload.RunAt)

		queueMutex.Lock()
		heap.Push(&jobQueue, job)
		jobQueueCond.Signal()
		queueMutex.Unlock()
	}
}
