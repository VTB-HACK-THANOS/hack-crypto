package api

import (
	"errors"
	"net/http"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
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
		switch err.(type) {
		case models.ForbiddenError:
			return ctx.String(http.StatusForbidden, "forbidden")
		case models.AlreadyExistsError:
			return ctx.String(http.StatusForbidden, "already exists")
		default:
			return ctx.String(http.StatusInternalServerError, "failed to register new user")
		}
	}

	return ctx.JSON(http.StatusOK, nil)
}

// handleUserBalance return user balance.
func (s *Server) handleUserBalance(ctx echo.Context) error {
	username, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if username == "" {
		echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	balance, err := s.UserManagementService.Balance(ctx.Request().Context(), username)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, balance)
}

// handleTransfer transfers currency to another user.
func (s *Server) handleTransfer(ctx echo.Context) error {
	type Request struct {
		ToUser string  `json:"to_user"`
		Amount float64 `json:"amount"`
	}
	var req Request

	if err := ctx.Bind(&req); err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	fromUser, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if fromUser == "" {
		echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	if err := s.UserManagementService.Transfer(ctx.Request().Context(), fromUser, req.ToUser, req.Amount); err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

// handleUserHistory returns user's transactions history
func (s *Server) handleUserHistory(ctx echo.Context) error {
	type Request struct {
		Page   int
		Offset int
		Sort   string
	}
	var r Request

	if err := ctx.Bind(&r); err != nil {
		echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if user == "" {
		echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	history, err := s.UserManagementService.History(ctx.Request().Context(), user, r.Page, r.Offset, r.Sort)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, history)
}
