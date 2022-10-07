package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anthonyaspen/emlvid-back/models"
)

// User returns a user by email.
func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{Email: email}
	if err := s.db.
		NewSelect().
		Model(user).
		Relation("Roles").
		WherePK().
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &models.NotFoundError{}
		}
		return nil, err
	}
	return user, nil
}

// UpsertUser inserts new user or updates if this user exists.
func (s *Storage) UpsertUser(ctx context.Context, user *models.User) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.NewInsert().
		Model(user).
		On("CONFLICT (email) DO UPDATE").
		Exec(ctx); err != nil {
		return err
	}

	if user.Roles != nil || len(user.Roles) != 0 {
		if _, err := tx.Exec("DELETE FROM user_roles WHERE user_email = ?", user.Email); err != nil {
			fmt.Println("here3")
			return err
		}

		userRoles := make([]*models.UserRole, 0, len(user.Roles))

		for _, role := range user.Roles {
			userRoles = append(userRoles, &models.UserRole{UserEmail: user.Email, RoleID: role.RoleID})
		}

		if _, err := tx.NewInsert().Model(&userRoles).Exec(ctx); err != nil {
			return err
		}
	}

	_ = tx.Commit()

	return err
}

// Roles returns list of existing roles.
func (s *Storage) Roles(ctx context.Context) ([]*models.Role, error) {
	roles := make([]*models.Role, 0)
	if err := s.db.NewSelect().Model(&roles).Scan(ctx); err != nil {
		return nil, err
	}

	return roles, nil
}
