package muxly_backend

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/Muxly-Corp/muxly-msg-subscriber/config"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/generated/muxly_backend"
)

type UserInfo struct {
	UserID   uuid.UUID
	Username string
}

type Client struct {
	inner *muxly_backend.ClientWithResponses
}

func NewClient(cfg config.MuxlyBackend) (*Client, error) {
	c, err := muxly_backend.NewClientWithResponses(cfg.URL,
		muxly_backend.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("X-Internal-API-Key", cfg.InternalAPIKey)
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("muxly_backend client: %w", err)
	}
	return &Client{inner: c}, nil
}

func (c *Client) ValidateToken(ctx context.Context, accessToken string) (*UserInfo, error) {
	resp, err := c.inner.ValidateTokenWithResponse(ctx, &muxly_backend.ValidateTokenParams{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("validate token: status %s", resp.Status())
	}
	return &UserInfo{
		UserID:   uuid.UUID(resp.JSON200.UserId),
		Username: resp.JSON200.Username,
	}, nil
}
