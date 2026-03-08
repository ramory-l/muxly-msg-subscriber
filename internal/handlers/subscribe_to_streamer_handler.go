package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Muxly-Corp/muxly-msg-subscriber/config"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/domains/platform"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/generated/msg_subscriber"
	"github.com/Muxly-Corp/muxly-shared/logger"
)

// Ensure sseResponse implements SubscribeToStreamerResponseObject.
var _ msg_subscriber.SubscribeToStreamerResponseObject = sseResponse{}

type SubscribeToStreamerHandler struct {
	service platform.Service
	cfg     config.Subscriber
}

func NewSubscribeToStreamerHandler(service platform.Service, cfg config.Subscriber) *SubscribeToStreamerHandler {
	return &SubscribeToStreamerHandler{service: service, cfg: cfg}
}

func (h *SubscribeToStreamerHandler) Handle(
	ctx context.Context,
	request msg_subscriber.SubscribeToStreamerRequestObject,
) (msg_subscriber.SubscribeToStreamerResponseObject, error) {
	streamerName := request.Params.To
	ch, unsubscribe, err := h.service.Subscribe(ctx, streamerName)
	if err != nil {
		logger.Errorf(ctx, "failed to subscribe to %s: %v", streamerName, err)

		return msg_subscriber.SubscribeToStreamer500JSONResponse{
			Error:   "subscribe_failed",
			Message: err.Error(),
		}, nil
	}

	return sseResponse{ch: ch, unsubscribe: unsubscribe, cfg: h.cfg}, nil
}

type sseResponse struct {
	ch          <-chan *platform.UnifiedMessage
	unsubscribe func()
	cfg         config.Subscriber
}

func (r sseResponse) VisitSubscribeToStreamerResponse(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "text/event-stream")
	c.Set(fiber.HeaderCacheControl, "no-cache")
	c.Set("Connection", "keep-alive")
	c.Status(fiber.StatusOK)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer r.unsubscribe()

		ticker := time.NewTicker(time.Duration(r.cfg.BatchIntervalMs) * time.Millisecond)
		defer ticker.Stop()

		batch := make(msg_subscriber.MessageBatch, 0, r.cfg.BatchMaxSize)

		flush := func() bool {
			if len(batch) == 0 {
				return true
			}
			data, err := json.Marshal(batch)
			if err != nil {
				batch = batch[:0]
				return true
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return false
			}
			if err := w.Flush(); err != nil {
				return false
			}
			batch = batch[:0]
			return true
		}

		for {
			select {
			case msg, ok := <-r.ch:
				if !ok {
					return
				}
				batch = append(batch, toGeneratedMessage(msg))
				if len(batch) >= r.cfg.BatchMaxSize {
					if !flush() {
						return
					}
				}
			case <-ticker.C:
				if !flush() {
					return
				}
			}
		}
	})

	return nil
}

func toGeneratedMessage(m *platform.UnifiedMessage) msg_subscriber.UnifiedMessage {
	return msg_subscriber.UnifiedMessage{
		Id:            m.ID,
		Platform:      m.Platform,
		Content:       m.Content,
		Timestamp:     m.Timestamp,
		IsAction:      m.IsAction,
		IsHighlighted: m.IsHighlighted,
		User:          toGeneratedUser(m.User),
		Emotes:        toGeneratedEmotes(m.Emotes),
	}
}

func toGeneratedUser(u platform.UnifiedUser) msg_subscriber.User {
	user := msg_subscriber.User{
		Id:            u.ID,
		Username:      u.Username,
		DisplayName:   u.DisplayName,
		IsModerator:   &u.IsModerator,
		IsSubscriber:  &u.IsSubscriber,
		IsBroadcaster: &u.IsBroadcaster,
	}
	if len(u.Badges) > 0 {
		user.Badges = &u.Badges
	}
	if u.Color != "" {
		user.Color = &u.Color
	}
	return user
}

func toGeneratedEmotes(emotes []platform.UnifiedEmote) *[]msg_subscriber.Emote {
	if len(emotes) == 0 {
		return nil
	}
	result := make([]msg_subscriber.Emote, len(emotes))
	for i, e := range emotes {
		result[i] = msg_subscriber.Emote{
			Id:       e.ID,
			Name:     e.Name,
			Url:      e.URL,
			Platform: e.Platform,
		}
	}
	return &result
}
