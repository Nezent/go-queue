package worker

import (
	"context"
	"log"

	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/Nezent/go-queue/internal/worker/enqueue"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartPgListener(ctx context.Context, channel string, pool *pgxpool.Pool, dispatcher *enqueue.TaskDispatcher, c *bootstrap.Container) {
	conn, _ := pool.Acquire(ctx)
	defer conn.Release()
	// Listen to the specified channel
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

		emailPayload, err := c.JobHandler.GetJobPayload(ctx, jobID)
		if err != nil {
			log.Printf("[LISTENER] Failed to fetch job payload for ID %s: %v\n", jobID, err)
			continue
		}
		sendJobEmail(ctx, *emailPayload, dispatcher)
		log.Printf("[LISTENER] Job email sent for ID %s\n", jobID)
	}
}

func sendJobEmail(context context.Context, payload task.EmailPayload, dispatcher *enqueue.TaskDispatcher) {

	_ = dispatcher.EnqueueSendJobEmail(context, task.EmailPayload{
		Recipient: payload.Recipient,
		Subject:   payload.Subject,
		Body:      payload.Body,
	})
}
