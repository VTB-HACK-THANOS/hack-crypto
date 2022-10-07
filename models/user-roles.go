package models

import "github.com/uptrace/bun"

type UserRole struct {
	bun.BaseModel `bun:"table:user_roles"`

	UserEmail string `bun:"user_email,pk" json:"-"`
	RoleID    uint   `bun:"role_id" json:"-"`
}
