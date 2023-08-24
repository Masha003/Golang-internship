package models

import "gorm.io/gorm"

// database model
type User struct {
	*gorm.Model
	Id       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
	Image    string `json: "path_to_image"`
}
