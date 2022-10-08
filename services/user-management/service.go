package users

import (
	"context"
	"errors"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/VTB-HACK-THANOS/hack-crypto/services/utils"
)

const (
	regularUserID = 1
)

type Service struct {
	store        Storage
	CryptoWallet CryptoWallet
}

func New(store Storage, cryptoWallet CryptoWallet) (*Service, error) {
	s := &Service{
		store:        store,
		CryptoWallet: cryptoWallet,
	}

	return s, nil
}

// InsertWhiteList inserts new user in a white list.
func (s *Service) InsertWhiteList(ctx context.Context, email string, createdBy string) error {
	if email == "" {
		return errors.New("email is empty")
	}

	return s.store.UpsertUser(ctx, &models.User{Email: email, ManagerEmail: createdBy})
}

// Roles returns list of roles.
func (s *Service) Roles(ctx context.Context) ([]*models.Role, error) {
	return s.store.Roles(ctx)
}

// RegisterUser register new user. Only if it's email already exists without a password.
func (s *Service) RegisterUser(ctx context.Context, email, password string) error {
	//return error if user has password already
	user, err := s.store.User(ctx, email)
	if errors.Is(err, &models.NotFoundError{}) {
		return &models.ForbiddenError{}
	}

	if user.PasswordHash != "" {
		return &models.AlreadyExistsError{}
	}

	user.PasswordHash, err = utils.HashAndSalt([]byte(password))
	if err != nil {
		return err
	}

	user.Roles = []*models.UserRole{
		{
			UserEmail: user.Email,
			RoleID:    regularUserID,
		},
	}

	wallet, err := s.CryptoWallet.CreateWallet()
	if err != nil {
		return err
	}

	//TODO transaction
	user.WalletCredentials.PrivateKey = wallet.PrivateKey
	user.WalletCredentials.PublicKey = wallet.PublicKey

	if err := s.store.UpsertUser(ctx, user); err != nil {
		return err
	}

	return nil
}

// Balance returns user's balance.
func (s *Service) Balance(ctx context.Context, username string) (*models.Balance, error) {
	u, err := s.store.User(ctx, username)
	if err != nil {
		return nil, err
	}

	balance, err := s.CryptoWallet.Balance(u.WalletCredentials.PublicKey)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (s *Service) Transfer(ctx context.Context, from, to string, amount float64) error {
	fromUser, err := s.store.User(ctx, from)
	if err != nil {
		return err
	}

	toUser, err := s.store.User(ctx, to)
	if err != nil {
		return err
	}

	if err := s.CryptoWallet.Transfer(
		fromUser.WalletCredentials.PrivateKey,
		toUser.WalletCredentials.PublicKey,
		amount,
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) History(ctx context.Context, username string, page, offset int, sort string) ([]*models.History, error) {
	u, err := s.store.User(ctx, username)
	if err != nil {
		return nil, err
	}

	history, err := s.CryptoWallet.History(ctx, u.PublicKey, page, offset, sort)
	if err != nil {
		return nil, err
	}

	return history, nil
}
