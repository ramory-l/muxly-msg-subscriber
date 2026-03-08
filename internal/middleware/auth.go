package middleware

import (
	muxlyctx "github.com/Muxly-Corp/muxly-shared/context"
	"github.com/gofiber/fiber/v2"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/generated/core"
	"github.com/Muxly-Corp/muxly-msg-subscriber/internal/generated/msg_subscriber"
	muxlyBackendClient "github.com/Muxly-Corp/muxly-msg-subscriber/internal/pkg/clients/muxly_backend"
)

type Auth struct {
	backendClient *muxlyBackendClient.Client
}

func NewAuthMiddleware(backendClient *muxlyBackendClient.Client) *Auth {
	return &Auth{backendClient: backendClient}
}

func (m *Auth) StrictMiddleware(f msg_subscriber.StrictHandlerFunc, operationID string) msg_subscriber.StrictHandlerFunc {
	return func(c *fiber.Ctx, request any) (any, error) {
		if _, requiresAuth := c.Context().UserValue(msg_subscriber.BearerAuthScopes).([]string); !requiresAuth {
			return f(c, request)
		}

		cookie := c.Cookies("access_token")
		if cookie == "" {
			return core.ErrorResponse{
				Error:   "unauthorized",
				Message: "Missing access_token cookie",
			}, nil
		}

		userInfo, err := m.backendClient.ValidateToken(c.UserContext(), cookie)
		if err != nil {
			return core.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired token",
			}, nil
		}

		ctx := c.UserContext()
		ctx = muxlyctx.WithUserID(ctx, userInfo.UserID)
		ctx = muxlyctx.WithUsername(ctx, userInfo.Username)
		ctx = muxlyctx.WithToken(ctx, cookie)
		c.SetUserContext(ctx)

		return f(c, request)
	}
}
