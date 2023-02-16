package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
	Token    string  `json:"token"`
	Icon     *string `json:"icon"`
}
