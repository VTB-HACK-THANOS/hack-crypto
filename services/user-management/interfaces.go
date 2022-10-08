package users

import (
	"context"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
)

type Storage interface {
	User(context.Context, string) (*models.User, error)
	UpsertUser(ctx context.Context, user *models.User) error
	Roles(ctx context.Context) ([]*models.Role, error)
}

type CryptoWallet interface {
	CreateWallet() (*models.WalletCredentials, error)
	Transfer(fromPrivateKey, toPublicKey string, amount float64) error
	Balance(publicKey string) (*models.Balance, error)
	History(ctx context.Context, publicKey string, page, offset int, sort string) ([]*models.History, error)
}
