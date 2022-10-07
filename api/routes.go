package api

import (
	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) setupRoutes(e *echo.Echo, rlCfg middleware.RateLimiterConfig) {
	e.Use(middleware.CORS())

	// api
	api := e.Group("/api")
	api.POST("/registration", s.handleRegistration, middleware.RateLimiterWithConfig(rlCfg))

	// v1
	v1 := api.Group("/v1", verifyAuth(s.AuthService))

	// admin endpoints
	admin := v1.Group("/admin", verifyUserAccessWrites(models.PrivilegedUser))
	// add new user in a whitelist
	admin.POST("/white-list", s.handleInsertWhiteList)

}
