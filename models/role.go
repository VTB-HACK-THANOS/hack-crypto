package models

import "github.com/uptrace/bun"

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID   uint   `bun:"id,pk" json:"id"`
	Name string `bun:"name" json:"name"`
}
