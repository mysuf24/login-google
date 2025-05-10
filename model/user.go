package model

import "gorm.io/gorm"

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"unique;not null"`
	Name        string
	Avatar      string
	CreatedAt   int64
	UpdatedAt   int64
	DeteletedAt gorm.DeletedAt `gorm:"index"`
}
