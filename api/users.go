package api

import (
	"errors"
	"fmt"
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

	createdBy, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if createdBy == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	if err := s.UserManagementService.InsertWhiteList(ctx.Request().Context(), req.Email, createdBy); err != nil {
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
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if username == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	balance, err := s.UserManagementService.Balance(ctx.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	fromUser, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if fromUser == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	if err := s.UserManagementService.Transfer(ctx.Request().Context(), fromUser, req.ToUser, req.Amount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if user == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	history, err := s.UserManagementService.History(ctx.Request().Context(), user, r.Page, r.Offset, r.Sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, history)
}

// handleTasks returns list of tasks
func (s *Server) handleTasks(ctx echo.Context) error {
	user, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}
	if user == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	resp, err := s.AssignmentService.Tasks(ctx.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) handleInsertTask(ctx echo.Context) error {
	var t models.Task

	if err := ctx.Bind(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, ok := ctx.Get(usernameCtx).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to convert to strig"))
	}

	if user == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("empty usernameCtx"))
	}

	fmt.Println(t)
	t.UserEmail = user

	resp, err := s.AssignmentService.InsertTask(ctx.Request().Context(), &t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Here")

	return ctx.JSON(http.StatusOK, resp)
}
