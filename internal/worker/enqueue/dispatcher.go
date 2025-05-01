package enqueue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/hibiken/asynq"
)

type TaskDispatcher struct {
	Client *asynq.Client
}

func NewTaskDispatcher(redisOpt asynq.RedisClientOpt) *TaskDispatcher {
	return &TaskDispatcher{
		Client: asynq.NewClient(redisOpt),
	}
}

func (d *TaskDispatcher) EnqueueSendVerificationEmail(ctx context.Context, payload task.SendVerificationEmailPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(task.TaskSendVerificationEmail, data)

	_, err = d.Client.EnqueueContext(ctx, task, asynq.MaxRetry(5), asynq.Timeout(30*time.Second))
	return err
}

func (d *TaskDispatcher) EnqueueSendJobEmail(ctx context.Context, payload task.EmailPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(task.TaskSendJobEmail, data)

	_, err = d.Client.EnqueueContext(ctx, task, asynq.MaxRetry(5), asynq.Timeout(30*time.Second))
	return err
}
