package bootstrap

import (
	"github.com/Nezent/go-queue/internal/handler"
	"github.com/Nezent/go-queue/internal/repository"
	"github.com/Nezent/go-queue/internal/service"
	"github.com/Nezent/go-queue/internal/websocket"
	"github.com/Nezent/go-queue/internal/worker/enqueue"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	UserHandler    handler.UserHandler
	JobHandler     handler.JobHandler
	TaskDispatcher *enqueue.TaskDispatcher
	WebSocketHub   *websocket.Hub
}

func Initialize(db *pgxpool.Pool, dispatcher *enqueue.TaskDispatcher, webSocketHub *websocket.Hub) *Container {

	return &Container{
		UserHandler: handler.UserHandler{
			Service: service.NewUserService(repository.NewUserRepository(db), dispatcher),
		},
		TaskDispatcher: dispatcher,
		JobHandler: handler.JobHandler{
			Service: service.NewJobService(repository.NewJobRepository(db)),
		},
		WebSocketHub: webSocketHub,
		// other handlers...
	}
}

func InitializeDispatcher(redisOpt asynq.RedisClientOpt) *enqueue.TaskDispatcher {
	return enqueue.NewTaskDispatcher(redisOpt)
}

func SetupWebSocketHub() *websocket.Hub {
	hub := websocket.NewHub()
	return hub
}
