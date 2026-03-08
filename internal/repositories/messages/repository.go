package messages

import (
	"encoding/json"
	"fmt"

	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/domains/platform"
	sharedqueue "github.com/Muxly-Corp/muxly-shared/queue"
)

type Repository struct{ queue sharedqueue.Queue }

func NewRepository(q sharedqueue.Queue) *Repository { return &Repository{queue: q} }

func (r *Repository) Subscribe(streamerName string, handler func(msg *platform.UnifiedMessage)) (func() error, error) {
	subject := "muxly.streamer." + streamerName
	unsubscribe, err := r.queue.Subscribe(subject, func(data []byte) {
		var msg platform.UnifiedMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}
		handler(&msg)
	})
	if err != nil {
		return nil, fmt.Errorf("messages: subscribe to %s: %w", streamerName, err)
	}
	return unsubscribe, nil
}
