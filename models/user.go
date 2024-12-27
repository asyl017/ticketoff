package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key;auto_increment"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
}

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(&User{})
}
