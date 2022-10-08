package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Email        string `bun:"email,pk" json:"email"`
	PasswordHash string `bun:"password" json:"-"`

	JobTitle string      `bun:"job_title" json:"job_title"`
	Name     string      `bun:"name" json:"name"`
	Roles    []*UserRole `bun:"roles,rel:has-many,join:email=user_email" json:"roles"`

	WalletCredentials
}

type Balance struct {
	RublesAmount float64 ` json:"coinsAmount"`
	MaticAmount  float64 `json:"maticAmount"`
}

type History struct {
	BlockNumber    float64   `json:"blockNumber"`
	TimeStamp      time.Time `json:"timeStamp"`
	ContactAddress string    `json:"contactAddress"`
	From           string    `json:"from"`
	To             string    `json:"to"`
	Value          float64   `json:"value"`
	TokenID        int       `json:"tokenId"`
	TokenName      string    `json:"tokenName"`
	TokenSymbol    string    `json:"tokenSymbol"`
	GasUsed        float64   `json:"gasUsed"`
	Confirmations  float64   `json:"confirmations"`
}
