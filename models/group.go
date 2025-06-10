package models

import "gorm.io/gorm"

// Group represents a group with a unique name
// and a many-to-many relationship with users.
type Group struct {
	gorm.Model
	Name  string  `json:"name" gorm:"unique;not null"`
	Users []*User `json:"users" gorm:"many2many:user_groups;"`
}
