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
	TasksByCreator(ctx context.Context, username string) ([]*models.Task, error)
	InsertTask(ctx context.Context, t *models.Task) (*models.Task, error)
	User(context.Context, string) (*models.User, error)
}
