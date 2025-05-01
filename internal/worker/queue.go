package worker

import (
	"github.com/Nezent/go-queue/internal/worker/processor"
	"github.com/Nezent/go-queue/internal/worker/task"
	"github.com/hibiken/asynq"
)

// Initializes and returns the Asynq server
func NewAsynqServer(redisOpt asynq.RedisClientOpt) *asynq.Server {
	return asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})
}

// Initializes and returns the task mux with handlers registered
func NewServeMux(processor *processor.TaskProcessor) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TaskSendVerificationEmail, processor.HandleSendVerificationEmail)
	mux.HandleFunc(task.TaskSendJobEmail, processor.HandleSendJobEmail)
	// Register other task handlers here
	return mux
}
