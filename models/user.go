package models

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`

	Email        string `bun:"email,pk" json:"email"`
	PasswordHash string `bun:"password" json:"-"`

	JobTitle string      `bun:"job_title" json:"job_title"`
	Name     string      `bun:"name" json:"name"`
	Roles    []*UserRole `bun:"roles,rel:has-many,join:email=user_email" json:"roles"`
}
