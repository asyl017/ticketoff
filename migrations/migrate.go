package migrations

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func InitDB(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to database", err)
		return nil, err
	}
	return db, nil
}
