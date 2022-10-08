package files

import (
	"context"
	"errors"
	"net/http"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/google/uuid"
)

type Service struct {
	Store Storage
}

func New(store Storage) (*Service, error) {
	s := &Service{Store: store}

	return s, nil
}

// Insert saves new file.
func (s *Service) Insert(ctx context.Context, file *models.Question) error {
	if len(file.Data) == 0 {
		return errors.New("data is empty")
	}

	if file.Name == "" {
		return errors.New("filename is empty")
	}

	if file.ContentType == "" {
		file.ContentType = http.DetectContentType(file.Data)
	}

	return s.Store.InsertQuestion(ctx, file)
}

// Preview returns list of audio without data.
func (s *Service) Preview(ctx context.Context, contentType string, limit, offset int) ([]*models.Question, error) {
	return nil, errors.New("not implemented")
}

// {
// 	ct := strings.TrimLeft(contentType, "/")
//
// 	files, err := s.Store.QuestionsByScopes(ctx, ct, limit, offset)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	result := make([]*models.Question, 0, len(files))
// 	for _, f := range files {
// 		result = append(result, &models.Question{ID: f.ID, Name: f.Name, ContentType: f.ContentType})
// 	}
//
// 	return result, nil
// }

// ByID returns audio by id with data.
func (s *Service) ByID(ctx context.Context, id uuid.UUID) (*models.Question, error) {
	if id == uuid.Nil {
		return nil, errors.New("not found")
	}

	return s.Store.QuestionByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}
