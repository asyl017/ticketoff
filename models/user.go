package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key;auto_increment"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(&User{})
}
