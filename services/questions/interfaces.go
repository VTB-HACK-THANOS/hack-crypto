package files

import (
	"context"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/google/uuid"
)

type Storage interface {
	InsertQuestion(ctx context.Context, q *models.Question) error
	QuestionByID(ctx context.Context, id uuid.UUID) (*models.Question, error)
	// QuestionsByScopes(ctx context.Context, contentType string, limit, offset int) ([]*models.Question, error)
	// DeleteFileByID(ctx context.Context, id uuid.UUID) error
}
