package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:" unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(&User{})
}
