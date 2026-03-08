package app

import (
	"github.com/gofiber/fiber/v2"
	sharedapp "github.com/Muxly-Corp/muxly-shared/app"
	sharedqueue "github.com/Muxly-Corp/muxly-shared/queue"
	"github.com/Muxly-Corp/muxly-msg-subscriber/config"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/generated/msg_subscriber"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/handlers"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/middleware"
)

func ProvideConfig(infra *sharedapp.Infra) *config.Config {
	return infra.Config.(*config.Config)
}

func ProvideNATSQueue(infra *sharedapp.Infra) *sharedqueue.NATSQueue {
	return infra.Queue.(*sharedqueue.NATSQueue)
}

// NewServiceSetup composes the ServiceSetup the framework needs to configure Fiber.
func NewServiceSetup(
	serviceRouter *handlers.ServiceHandlerRouter,
	authMiddleware *middleware.Auth,
) *sharedapp.ServiceSetup {
	return &sharedapp.ServiceSetup{
		Router: func(f *fiber.App) {
			serviceServer := msg_subscriber.NewStrictHandler(
				serviceRouter,
				[]msg_subscriber.StrictMiddlewareFunc{authMiddleware.StrictMiddleware},
			)
			msg_subscriber.RegisterHandlers(f, serviceServer)
		},
	}
}
