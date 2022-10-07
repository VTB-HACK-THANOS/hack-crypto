package auth

import (
	"context"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
)

type Storage interface {
	User(context.Context, string) (*models.User, error)
}
