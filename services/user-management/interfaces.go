package users

import (
	"context"

	"github.com/anthonyaspen/emlvid-back/models"
)

type Storage interface {
	User(context.Context, string) (*models.User, error)
	UpsertUser(ctx context.Context, user *models.User) error
	Roles(ctx context.Context) ([]*models.Role, error)
}
