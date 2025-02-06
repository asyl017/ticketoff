package models

import (
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email            string             `json:"email" bson:"email"`
	Password         string             `json:"password" bson:"password"`
	EmailConfirmed   bool               `json:"email_confirmed" bson:"email_confirmed"`
	VerificationCode string             `json:"verification_code" bson:"verification_code"`
	Role             string             `json:"role" bson:"role"`
	IsAdmin          bool               `json:"is_admin" bson:"is_admin"`
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
