package worker

import (
	"container/heap"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Nezent/go-queue/internal/bootstrap"
	"github.com/Nezent/go-queue/internal/worker/enqueue"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/google/uuid"
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

func InitJobQueue(ctx context.Context, dispatcher *enqueue.TaskDispatcher, c *bootstrap.Container) {
	heap.Init(&jobQueue)
	go processJobs(ctx, &jobQueue, dispatcher, c)
}

func processJobs(ctx context.Context, jobQueue *JobPriorityQueue, dispatcher *enqueue.TaskDispatcher, c *bootstrap.Container) {
	for {
		queueMutex.Lock()

		for len(*jobQueue) == 0 {
			// Wait until a job is pushed
			jobQueueCond.Wait()
		}

		jobItem := heap.Pop(jobQueue).(*JobItem)

		now := time.Now()
		if now.Before(jobItem.RunAt) {
			// Requeue and wait until the job is ready
			heap.Push(jobQueue, jobItem)
			sleepDuration := jobItem.RunAt.Sub(now)
			queueMutex.Unlock()

			time.Sleep(sleepDuration) // wait until the job's RunAt time
			continue
		}

		queueMutex.Unlock()

		// Process the job
		err := dispatcher.EnqueueSendJobEmail(ctx, jobItem.Payload)
		if err != nil {
			// Requeue on error
			jobItem.Attempts++
			log.Printf("[PROCESS] Job ID %s failed, retrying... (Attempt %d)\n", jobItem.ID, jobItem.Attempts)

			queueMutex.Lock()
			heap.Push(jobQueue, jobItem)
			jobQueueCond.Signal()
			queueMutex.Unlock()
		} else {
			jsonMsg := task.WebSocketPayload{
				JobID:   jobItem.ID.String(),
				JobType: jobItem.JobType,
				Status:  jobItem.Status,
			}
			jsonMsgBytes, err := json.Marshal(jsonMsg)
			if err != nil {
				log.Printf("[LISTENER] Failed to marshal JSON for job ID %s: %v\n", jobItem.ID, err)
				continue
			}
			c.WebSocketHub.Broadcast <- jsonMsgBytes
			log.Printf("[PROCESS] Job ID %s executed successfully.\n", jobItem.ID)
		}
	}
}
