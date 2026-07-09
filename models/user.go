package models

import "gorm.io/gorm"


type User struct{
	gorm.Model
	Name string 
	Email string
	Password string
}
type RegistrationInput struct{
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32,excludes= "`
}
type LoginInput struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32,excludes= "`
}