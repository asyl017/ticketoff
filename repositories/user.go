package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import Postgres dialect
	"log"
	"ticketoff/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(updatedUser *models.User) (*models.User, error)
	DeleteUser(id string) error
	GetUsers() ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepository) CreateUser(user *models.User) error {
	if u.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// Log the incoming user data for debugging
	log.Printf("Creating user: %+v\n", user)

	// Attempt to create the user in the database
	if err := u.db.Create(user).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		return err
	}

	return nil
}

func (u userRepository) GetUserByID(id string) (*models.User, error) {
	if u.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var user models.User
	log.Printf("Fetching user with ID: %s\n", id)

	if err := u.db.First(&user, "id = ?", id).Error; err != nil {
		log.Printf("Error fetching user: %v\n", err)
		return nil, err
	}

	log.Printf("Fetched user: %+v\n", user)
	return &user, nil
}

// UpdateUser updates a user's details
func (u userRepository) UpdateUser(updatedUser *models.User) (*models.User, error) {
	model := updatedUser
	err := u.db.Save(updatedUser).Error
	if err != nil {
		log.Printf("Error updating user: %v\n", err)
		return nil, err
	}
	return model, err
}

// DeleteUser removes a user from the database by ID
func (u userRepository) DeleteUser(id string) error {
	var user models.User

	// Find user by ID
	if err := u.db.First(&user, "id = ?", id).Error; err != nil {
		return err
	}

	// Delete user
	return u.db.Delete(&user).Error
}

// GetUsers retrieves all users from the database
func (u userRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := u.db.Find(&users).Error
	return users, err
}
