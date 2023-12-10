package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required, min=5"`
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=6, oneof=uppercase&lowercase&numeric"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=5"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}
