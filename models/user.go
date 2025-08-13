package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"size:191"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Followers []*User `gorm:"many2many:user_followers;joinForeignKey:UserID;joinReferences:FollowerID"`
	Following []*User `gorm:"many2many:user_following;joinForeignKey:UserID;joinReferences:FollowingID"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=5"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=5"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (createDto *CreateUserRequest) ToUser() *User {
	return &User{Username: createDto.Username, Email: createDto.Email, Password: createDto.Password}
}
