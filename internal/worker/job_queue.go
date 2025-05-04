package worker

import (
	"container/heap"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Nezent/go-queue/common"
	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/Nezent/go-queue/internal/middleware"
	"github.com/Nezent/go-queue/internal/worker/enqueue"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobItem struct {
	ID       uuid.UUID
	RunAt    time.Time
	Priority int
	Attempts int
	Payload  task.EmailPayload
	JobType  string
	Status   string
	index    int // required by heap.Interface
}

type JobPriorityQueue []*JobItem

func (pq JobPriorityQueue) Len() int { return len(pq) }

func (pq JobPriorityQueue) Less(i, j int) bool {
	if pq[i].RunAt.Equal(pq[j].RunAt) {
		return pq[i].Priority < pq[j].Priority // lower value = higher priority
	}
	return pq[i].RunAt.Before(pq[j].RunAt)
}

func (pq JobPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *JobPriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*JobItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *JobPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

var (
	jobQueue     JobPriorityQueue
	queueMutex   sync.Mutex
	jobQueueCond = sync.NewCond(&queueMutex)
)

func InitJobQueue(ctx context.Context, dispatcher *enqueue.TaskDispatcher, c *bootstrap.Container, db *pgxpool.Pool) {
	heap.Init(&jobQueue)
	go processJobs(ctx, &jobQueue, dispatcher, c, db)
}

func processJobs(ctx context.Context, jobQueue *JobPriorityQueue, dispatcher *enqueue.TaskDispatcher, c *bootstrap.Container, db *pgxpool.Pool) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			queueMutex.Lock()

			if len(*jobQueue) == 0 {
				log.Println("Job queue is empty, waiting for jobs...")
				jobQueueCond.Wait()

				if ctx.Err() != nil {
					queueMutex.Unlock()
					return
				}
				queueMutex.Unlock()
				continue
			}

			now := time.Now().In(common.DhakaTZ).Truncate(time.Second)
			nextJob := (*jobQueue)[0]

			if nextJob.RunAt.After(now) {
				queueMutex.Unlock()
				continue
			}
			// Remove the job from the queue
			heap.Pop(jobQueue)
			queueMutex.Unlock()
			// Process the job
			log.Printf("[PROCESS] Processing job ID %s...\n", nextJob.ID)

			nextJob.Attempts++
			if nextJob.Attempts > 3 {
				log.Printf("[PROCESS] Job ID %s failed after 3 attempts, marking as failed.\n", nextJob.ID)
				nextJob.Status = "failed"
				err := updateJobStatus(ctx, nextJob.ID, nextJob.Status, nextJob.Attempts, c, db)
				if err != nil {
					log.Printf("[ERROR] Failed to update job status: %v\n", err)
				}
				continue
			}

			err := dispatcher.EnqueueSendJobEmail(ctx, nextJob.Payload)
			if err != nil {
				nextJob.Attempts++
				log.Printf("[PROCESS] Job ID %s failed, retrying... (Attempt %d)\n", nextJob.ID, nextJob.Attempts)
				nextJob.Status = "processing"
				nextJob.RunAt = time.Now().Add(time.Duration(nextJob.Attempts) * time.Minute)
				nextJob.Priority = 1
				err = updateJobStatus(ctx, nextJob.ID, nextJob.Status, nextJob.Attempts, c, db)
				if err != nil {
					log.Printf("[ERROR] Failed to update job status: %v\n", err)
				}
				queueMutex.Lock()
				heap.Push(jobQueue, nextJob)
				jobQueueCond.Signal()
				queueMutex.Unlock()
			} else {
				jsonMsg := task.WebSocketPayload{
					JobID:   nextJob.ID.String(),
					JobType: nextJob.JobType,
					Status:  "completed",
				}
				nextJob.Status = "completed"
				err = updateJobStatus(ctx, nextJob.ID, nextJob.Status, nextJob.Attempts, c, db)
				if err != nil {
					log.Printf("[ERROR] Failed to update job status: %v\n", err)
				}
				jsonMsgBytes, err := json.Marshal(jsonMsg)
				if err != nil {
					log.Printf("[LISTENER] Failed to marshal JSON for job ID %s: %v\n", nextJob.ID, err)
					continue
				}
				c.WebSocketHub.Broadcast <- jsonMsgBytes
				log.Printf("[PROCESS] Job ID %s executed successfully.\n", nextJob.ID)
			}
		}
	}
}

func updateJobStatus(ctx context.Context, jobID uuid.UUID, status string, attempts int, c *bootstrap.Container, db *pgxpool.Pool) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("failed to start transaction:", err)
		return err
	}
	defer tx.Rollback(ctx)

	ctx = context.WithValue(ctx, middleware.TxKey, tx)
	appErr := c.JobHandler.UpdateJobStatus(ctx, jobID, status, attempts)
	if appErr != nil {
		return appErr
	}
	return tx.Commit(ctx)
}
