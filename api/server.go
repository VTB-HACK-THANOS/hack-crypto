package api

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type Server struct {
	e                     *echo.Echo
	log                   Logger
	AuthService           AuthService
	UserManagementService UserManagementService
	rlCfg                 middleware.RateLimiterConfig
}

const (
	oneRequestPerSecond = 1
)

func New(
	authService AuthService,
	userManagementService UserManagementService,
	rlRate, rlExpiresIn time.Duration,
	rlBurst int,
) *Server {
	e := echo.New()
	e.HideBanner = false
	e.HidePort = false

	s := &Server{
		e:                     e,
		AuthService:           authService,
		UserManagementService: userManagementService,
	}

	s.setupRoutes(e, setupRateLimiter(rlRate, rlExpiresIn, rlBurst))

	return s
}

func (s *Server) Start(ctx context.Context, addr string) error {
	go func() error {
		select {
		case <-ctx.Done():
			if err := s.e.Shutdown(ctx); err != nil {
				return err
			}
		}
		return nil
	}()

	return s.e.Start(addr)
}

// setupRateLimiter returns rate limiter config with given configurations.
func setupRateLimiter(rlRate, rlExpiresIn time.Duration, rlBurst int) middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Every(rlRate),
				Burst:     rlBurst,
				ExpiresIn: rlExpiresIn,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
}
