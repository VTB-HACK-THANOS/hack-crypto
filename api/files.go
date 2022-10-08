package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	formFileName = "file"
)

// handleUploadFile handles uploading new audio background track.
func (s *Server) handleUploadFile(ctx echo.Context) error {
	form, err := ctx.FormFile(formFileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	f, err := form.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fileData, err := io.ReadAll(f)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println(len(fileData))

	file := &models.Question{
		ID:   uuid.New(),
		Name: form.Filename,
		Data: fileData,
	}

	if err := s.QuestionService.Insert(ctx.Request().Context(), file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusCreated)
}

// handleQuestionList returns list of files previews.
func (s *Server) handleQuestionList(ctx echo.Context) error {
	type Request struct {
		Limit       int    `query:"limit"`
		Offset      int    `query:"offset"`
		ContentType string `query:"content-type"`
	}
	var r Request
	if err := ctx.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tracks, err := s.QuestionService.Preview(ctx.Request().Context(), r.ContentType, r.Limit, r.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, tracks)
}

// handleQuestionByID returns file by id.
func (s *Server) handleQuestionByID(ctx echo.Context) error {
	type Request struct {
		ID uuid.UUID `param:"id"`
	}
	var r Request

	if err := ctx.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if r.ID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "uuid is nil or empty")
	}

	track, err := s.QuestionService.ByID(ctx.Request().Context(), r.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Blob(http.StatusOK, track.ContentType, track.Data)
}

// handleDeleteFileByID deletes file by id.
func (s *Server) handleDeleteFileByID(ctx echo.Context) error {
	type Request struct {
		ID uuid.UUID `param:"id"`
	}
	var r Request

	if err := ctx.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if r.ID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "uuid is nil or empty")
	}

	if err := s.QuestionService.Delete(ctx.Request().Context(), r.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
