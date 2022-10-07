package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// handleInsertWhiteList adds new user in a white list.
func (s *Server) handleInsertWhiteList(ctx echo.Context) error {
	type Req struct {
		Email string `json:"email"`
	}
	var req Req

	if err := ctx.Bind(&req); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if req.Email == "" {
		return ctx.String(http.StatusBadRequest, "email is empty")
	}

	if err := s.UserManagementService.InsertWhiteList(ctx.Request().Context(), req.Email); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.NoContent(http.StatusCreated)
}

// handleListRoles returns list of roles.
func (s *Server) handleListRoles(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, nil)
}

// handleRegistration registers a user.
func (s *Server) handleRegistration(ctx echo.Context) error {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req

	if err := ctx.Bind(&req); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if req.Email == "" {
		return ctx.String(http.StatusBadRequest, "email is empty")
	}

	if req.Password == "" {
		return ctx.String(http.StatusBadRequest, "password is empty")
	}

	if err := s.UserManagementService.RegisterUser(ctx.Request().Context(), req.Email, req.Password); err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to register new user")
	}

	return ctx.JSON(http.StatusOK, nil)
}
