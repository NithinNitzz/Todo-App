package models

import _ "github.com/go-playground/validator/v10"

type Users struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username" validate:"required,min=3,max=15"`
	Password string `json:"password" validate:"required,min=3,max=15"`
}
