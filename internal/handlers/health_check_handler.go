package handlers

import (
	"context"
	"time"

	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/generated/msg_subscriber"
)

type HealthCheckHandler struct{}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Handle(
	ctx context.Context,
	request msg_subscriber.HealthCheckRequestObject,
) (msg_subscriber.HealthCheckResponseObject, error) {
	return msg_subscriber.HealthCheck200JSONResponse{
		Status:    "ok",
		Timestamp: time.Now().UnixMilli(),
	}, nil
}
