package models

import (
	"github.com/jinzhu/gorm"
)

type Film struct {
	ID          uint   `json:"id" gorm:"primary_key;auto_increment"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
	ReleaseYear int    `json:"release_year"`
}

func (Film) TableName() string {
	return "movies"
}

func MigrateFilm(DB *gorm.DB) {
	DB.AutoMigrate(&Film{})
}
