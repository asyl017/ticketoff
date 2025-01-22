package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	ID             uint   `json:"id" gorm:"primary_key;auto_increment"`
	Email          string `json:"email" gorm:"unique"`
	Password       string `json:"password" gorm:"not null"`
	EmailConfirmed bool   `json:"email_confirmed" gorm:"default:false"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Role struct {
	ID          uint         `json:"id" gorm:"primary_key"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

func MigrateUser(DB *gorm.DB) {
	DB.AutoMigrate(&User{})
}
