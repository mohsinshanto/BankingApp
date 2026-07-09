package models

import "gorm.io/gorm"

type Account struct{
	gorm.Model
	AccountNo string   `gorm:"type:varchar(20);uniqueIndex;not null"`
	Balance   float64  `gorm:"not null;default:0"`
	Status    string   `gorm:"type:varchar(20);not null;default:'ACTIVE'"`
	UserID    uint     `gorm:"uniqueIndex;not null"`
	User      User     `gorm:"foreignKey:UserID"`
}