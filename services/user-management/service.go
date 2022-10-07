package users

import (
	"context"
	"errors"

	"github.com/VTB-HACK-THANOS/hack-crypto/models"
	"github.com/VTB-HACK-THANOS/hack-crypto/services/utils"
)

const (
	regularUserID = 2
)

type Service struct {
	store Storage
}

func New(store Storage) (*Service, error) {
	s := &Service{store: store}

	return s, nil
}

// InsertWhiteList inserts new user in a white list.
func (s *Service) InsertWhiteList(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("email is empty")
	}

	return s.store.UpsertUser(ctx, &models.User{Email: email})
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

	if err := s.store.UpsertUser(ctx, user); err != nil {
		return err
	}

	return nil
}
