package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticketoff/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(updatedUser *models.User) (*models.User, error)
	DeleteUser(id string) error
	GetUsers() ([]models.User, error)
	ConfirmEmail(email string) error
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) CreateUser(user *models.User) error {
	collection := u.db.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

func (u userRepository) GetUserByEmail(email string) (*models.User, error) {
	collection := u.db.Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (u userRepository) GetUserByID(id string) (*models.User, error) {
	collection := u.db.Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&user)
	return &user, err
}

func (u userRepository) UpdateUser(updatedUser *models.User) (*models.User, error) {
	collection := u.db.Collection("users")
	_, err := collection.UpdateOne(context.Background(), bson.M{"id": updatedUser.ID}, bson.M{"$set": updatedUser})
	return updatedUser, err
}

func (u userRepository) DeleteUser(id string) error {
	collection := u.db.Collection("users")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	return err
}

func (u userRepository) GetUsers() ([]models.User, error) {
	collection := u.db.Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u userRepository) ConfirmEmail(email string) error {
	collection := u.db.Collection("users")
	_, err := collection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"email_confirmed": true}})
	return err
}
