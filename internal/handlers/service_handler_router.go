package handlers

import (
	"context"

	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/generated/msg_subscriber"
)

type ServiceHandlerRouter struct {
	healthCheckHandler        *HealthCheckHandler
	subscribeToStreamerHandler *SubscribeToStreamerHandler
}

func NewServiceHandlerRouter(
	healthCheckHandler *HealthCheckHandler,
	subscribeToStreamerHandler *SubscribeToStreamerHandler,
) *ServiceHandlerRouter {
	return &ServiceHandlerRouter{
		healthCheckHandler:        healthCheckHandler,
		subscribeToStreamerHandler: subscribeToStreamerHandler,
	}
}

func (r *ServiceHandlerRouter) HealthCheck(ctx context.Context, request msg_subscriber.HealthCheckRequestObject) (msg_subscriber.HealthCheckResponseObject, error) {
	return r.healthCheckHandler.Handle(ctx, request)
}

func (r *ServiceHandlerRouter) SubscribeToStreamer(ctx context.Context, request msg_subscriber.SubscribeToStreamerRequestObject) (msg_subscriber.SubscribeToStreamerResponseObject, error) {
	return r.subscribeToStreamerHandler.Handle(ctx, request)
}

var _ msg_subscriber.StrictServerInterface = (*ServiceHandlerRouter)(nil)
