package models

import "gorm.io/gorm"

type Transaction struct{
	gorm.Model
	FromAccount string  `gorm:"type:varchar(20)"`
	ToAccount   string  `gorm:"type:varchar(20)"`
	Amount      float64 `gorm:"not null"`
	Type        string  `gorm:"type:varchar(20);not null"`
	Description string  `gorm:"type:varchar(255)"`
	Status      string  `gorm:"type:varchar(20);default:'SUCCESS'"`
}