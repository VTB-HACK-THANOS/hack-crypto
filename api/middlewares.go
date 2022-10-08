package api

import (
	"log"
	"net/http"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/labstack/echo/v4"
)

const (
	userAccessWritesCtx = "user_access_writes"
	usernameCtx         = "username_ctx"
)

// verifyToken verifies token by trying verifyAuthorization() and verifyCookie() methods.
func verifyAuth(auth AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			login, pass, ok := ctx.Request().BasicAuth()
			if !ok {
				return ctx.String(http.StatusBadRequest, "auth data is empty")
			}
			if login == "" {
				return ctx.String(http.StatusBadRequest, "login is empty")
			}
			if pass == "" {
				return ctx.String(http.StatusBadRequest, "pass is empty")
			}

			userAccessWrites, err := auth.AuthUser(ctx.Request().Context(), login, pass)
			if userAccessWrites == 0 || err != nil {
				return ctx.String(http.StatusForbidden, "failed to authorize user")
			}

			ctx.Set(userAccessWritesCtx, userAccessWrites)
			ctx.Set(usernameCtx, login)

			return next(ctx)
		}
	}
}

// verifyUserAccessWrites compares given and presented access writes
func verifyUserAccessWrites(requiredAccessWrites models.AccessWrites) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			userAccessWrites, ok := ctx.Get(userAccessWritesCtx).(models.AccessWrites)
			if !ok {
				return ctx.String(http.StatusForbidden, "failed to get user access writes")
			}

			if requiredAccessWrites > userAccessWrites {
				return ctx.NoContent(http.StatusForbidden)
			}
			log.Printf("req: %v, presented: %v\n", requiredAccessWrites, userAccessWrites)

			return next(ctx)
		}
	}
}
