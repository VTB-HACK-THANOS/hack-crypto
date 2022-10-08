package api

import (
	"context"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/google/uuid"
)

type AuthService interface {
	AuthUser(ctx context.Context, username, password string) (models.AccessWrites, error)
}

type UserManagementService interface {
	InsertWhiteList(ctx context.Context, email string) error
	Roles(ctx context.Context) ([]*models.Role, error)
	RegisterUser(context.Context, string, string) error
	Balance(context.Context, string) (*models.Balance, error)
	Transfer(ctx context.Context, from, to string, amount float64) error
	History(ctx context.Context, username string, page, offset int, sort string) ([]*models.History, error)
}

type QuestionService interface {
	Insert(context.Context, *models.Question) error
	Preview(ctx context.Context, contentType string, limit, offset int) ([]*models.Question, error)
	ByID(ctx context.Context, id uuid.UUID) (*models.Question, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Logger interface {
}
