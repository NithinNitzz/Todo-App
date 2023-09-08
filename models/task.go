package models

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Task struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name"  example:"learn golang"`
	Desc      string    `json:"description"  `
	DueDate   time.Time `json:"due_date" `
	Status    string    `json:"status"  `
	CreatedAT time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	UserId    int       `json:"user_id" `
	User      Users     `json:"user"  `
}
