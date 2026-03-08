package platform

import (
	"context"
	"sync"

	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/domains/platform"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/repositories/messages"
)

type Service struct {
	msgRepo   *messages.Repository
	mu        sync.Mutex
	listeners map[string][]chan *platform.UnifiedMessage
	natsSubs  map[string]func() error
}

func NewService(msgRepo *messages.Repository) *Service {
	return &Service{
		msgRepo:   msgRepo,
		listeners: make(map[string][]chan *platform.UnifiedMessage),
		natsSubs:  make(map[string]func() error),
	}
}

func (s *Service) Subscribe(_ context.Context, streamerName string) (<-chan *platform.UnifiedMessage, func(), error) {
	ch := make(chan *platform.UnifiedMessage, 16)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.natsSubs[streamerName]; !ok {
		unsubNATS, err := s.msgRepo.Subscribe(streamerName, func(msg *platform.UnifiedMessage) {
			s.broadcast(streamerName, msg)
		})
		if err != nil {
			return nil, nil, err
		}
		s.natsSubs[streamerName] = unsubNATS
	}

	s.listeners[streamerName] = append(s.listeners[streamerName], ch)

	unsubscribe := func() {
		s.mu.Lock()

		listeners := s.listeners[streamerName]
		for i, c := range listeners {
			if c == ch {
				s.listeners[streamerName] = append(listeners[:i], listeners[i+1:]...)
				close(ch)
				break
			}
		}

		var natsUnsub func() error
		if len(s.listeners[streamerName]) == 0 {
			if unsub, ok := s.natsSubs[streamerName]; ok {
				natsUnsub = unsub
				delete(s.natsSubs, streamerName)
			}
			delete(s.listeners, streamerName)
		}

		s.mu.Unlock()

		if natsUnsub != nil {
			_ = natsUnsub()
		}
	}

	return ch, unsubscribe, nil
}

func (s *Service) broadcast(streamerName string, msg *platform.UnifiedMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ch := range s.listeners[streamerName] {
		select {
		case ch <- msg:
		default: // drop message for slow consumers
		}
	}
}
