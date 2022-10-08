package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TaskType string

const (
	EasyTask   TaskType = "easy"
	MediumTask TaskType = "medium"
	HardTask   TaskType = "hard"
)

type Task struct {
	bun.BaseModel `bun:"table:tasks"`
	ID            uuid.UUID `bun:"id" json:"id"`
	Name          string    `bun:"name" json:"name"`
	Description   string    `bun:"description" json:"description"`
	UserEmail     string    `bun:"user_email" json:"user_email"`
	Type          string    `bun:"type" json:"type"`
}
