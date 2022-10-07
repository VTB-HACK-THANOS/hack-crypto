package auth

import (
	"context"

	"github.com/anthonyaspen/emlvid-back/models"
)

type Storage interface {
	User(context.Context, string) (*models.User, error)
}
