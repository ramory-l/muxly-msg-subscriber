package platform

import (
	"time"
)

type UnifiedMessage struct {
	ID            string
	Platform      string
	User          UnifiedUser
	Content       string
	Emotes        []UnifiedEmote
	Timestamp     time.Time
	IsAction      bool
	IsHighlighted bool
	TargetRoom    string // streamer username; repository publishes to muxly.streamer.<TargetRoom>
}

type UnifiedUser struct {
	ID            string
	Username      string
	DisplayName   string
	Badges        []string
	Color         string
	IsModerator   bool
	IsSubscriber  bool
	IsBroadcaster bool
}

type UnifiedEmote struct {
	ID       string
	Name     string
	URL      string
	Platform string
}
