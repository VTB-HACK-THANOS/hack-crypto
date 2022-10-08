package storage

import (
	"context"
	"errors"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/google/uuid"
)

// InsertFile saves new file.
func (s *Storage) InsertQuestion(ctx context.Context, q *models.Question) error {
	if _, err := s.db.
		NewInsert().
		Model(q).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Storage) QuestionByID(ctx context.Context, id uuid.UUID) (*models.Question, error) {
	return nil, errors.New("not implemented")
}
