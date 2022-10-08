package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Nft struct {
	bun.BaseModel `bun:"table:nfts"`

	ID        uuid.UUID `bun:"id" json:"id"`
	UserEmail string    `bun:"user_email" json:"user_email"`
	Type      string    `bun:"type" json:"type"`
}
