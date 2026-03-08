package platform

import "context"

type Service interface {
	Subscribe(ctx context.Context, streamerName string) (<-chan *UnifiedMessage, func(), error)
}
