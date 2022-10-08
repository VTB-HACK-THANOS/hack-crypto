package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Nft struct {
	bun.BaseModel `bun:"table:nfts"`

	ID          uuid.UUID `bun:"id" json:"id"`
	Text        string    `bun:"text" json:"text"`
	Name        string    `bun:"name" json:"name"`
	ContentType string    `bun:"content_type" json:"content_type"`
	Data        []byte    `bun:"data" json:"data,omitempty"`
}
