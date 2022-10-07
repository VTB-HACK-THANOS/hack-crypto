package auth

import (
	"context"
	"fmt"

	"github.com/anthonyaspen/emlvid-back/models"
	"github.com/anthonyaspen/emlvid-back/services/utils"
)

type Service struct {
	store Storage
}

func New(store Storage) (*Service, error) {
	s := &Service{store: store}

	return s, nil
}

func (svc *Service) AuthUser(ctx context.Context, username, password string) (models.AccessWrites, error) {
	u, err := svc.store.User(ctx, username)
	if err != nil {
		return models.NotAuthorizedUser, err
	}

	if u.PasswordHash == "" {
		return models.NotAuthorizedUser, &models.NotFoundError{}
	}

	isOk, err := utils.ComparePasswords(u.PasswordHash, password)
	if err != nil {
		return models.NotAuthorizedUser, err
	}
	if !isOk {
		return models.NotAuthorizedUser, fmt.Errorf("failed to authorizate a user %s", username)
	}

	var userAccessWrites models.AccessWrites

	for _, role := range u.Roles {
		if role.RoleID > uint(userAccessWrites) {
			userAccessWrites = models.AccessWrites(role.RoleID)
		}
	}

	return userAccessWrites, nil
}
