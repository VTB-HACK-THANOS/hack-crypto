package storage

import (
	"context"
	"database/sql"
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

func (s *Storage) TasksByCreator(ctx context.Context, username string) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)
	if err := s.db.
		NewSelect().
		Model(&tasks).
		Where("user_email = ?", username).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &models.NotFoundError{}
		}
		return nil, err
	}
	return tasks, nil
}

func (s *Storage) InsertTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	if _, err := s.db.
		NewInsert().
		Model(t).
		Returning("*").
		Exec(ctx); err != nil {
		return nil, err
	}

	return t, nil
}
