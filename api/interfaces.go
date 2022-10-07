package api

import (
	"context"

	"github.com/anthonyaspen/emlvid-back/models"
)

type AuthService interface {
	AuthUser(ctx context.Context, username, password string) (models.AccessWrites, error)
}

type UserManagementService interface {
	InsertWhiteList(ctx context.Context, email string) error
	Roles(ctx context.Context) ([]*models.Role, error)
	RegisterUser(context.Context, string, string) error
}

type Logger interface {
}
