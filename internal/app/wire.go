//go:build wireinject

package app

import (
	"github.com/google/wire"
	sharedapp "github.com/Muxly-Corp/muxly-shared/app"
	sharedqueue "github.com/Muxly-Corp/muxly-shared/queue"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/domains/platform"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/handlers"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/middleware"
	muxlyBackendClient "github.com/Muxly-Corp/muxly-msg-subscriber/internal/pkg/clients/muxly_backend"
	messagesRepo "github.com/Muxly-Corp/muxly-msg-subscriber/internal/repositories/messages"
	platformService "github.com/Muxly-Corp/muxly-msg-subscriber/internal/services/platform"
	"github.com/Muxly-Corp/muxly-msg-subscriber/config"
)

var ConfigSet = wire.NewSet(ProvideConfig)

var QueueSet = wire.NewSet(
	ProvideNATSQueue,
	wire.Bind(new(sharedqueue.Queue), new(*sharedqueue.NATSQueue)),
)

var NATSRepoSet = wire.NewSet(messagesRepo.NewRepository)

var MuxlyBackendSet = wire.NewSet(
	wire.FieldsOf(new(*config.Config), "MuxlyBackend"),
	muxlyBackendClient.NewClient,
	middleware.NewAuthMiddleware,
)

var PlatformServiceSet = wire.NewSet(
	platformService.NewService,
	wire.Bind(new(platform.Service), new(*platformService.Service)),
)

var HandlerSet = wire.NewSet(
	wire.FieldsOf(new(*config.Config), "Subscriber"),
	handlers.NewHealthCheckHandler,
	handlers.NewSubscribeToStreamerHandler,
	handlers.NewServiceHandlerRouter,
)

var AppSet = wire.NewSet(NewServiceSetup)

func InitializeApp(infra *sharedapp.Infra) (*sharedapp.ServiceSetup, error) {
	wire.Build(ConfigSet, QueueSet, NATSRepoSet, MuxlyBackendSet, PlatformServiceSet, HandlerSet, AppSet)
	return nil, nil
}
