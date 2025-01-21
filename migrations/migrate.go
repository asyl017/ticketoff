package migrations

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"ticketoff/models"
)

var db *gorm.DB

// InitDB initializes the database and runs migrations
func InitDB(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema
	if err := db.AutoMigrate(&models.User{}).Error; err != nil {
		log.Fatalf("Migration failed: %v", err)
		return nil, err
	}
	// Automatically migrate the Film schema
	if err := db.AutoMigrate(&models.Film{}).Error; err != nil {
		log.Fatalf("Film migration failed: %v", err)
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil

}

// GetDB returns the DB instance for use in other packages
func GetDB() *gorm.DB {
	return db
}
