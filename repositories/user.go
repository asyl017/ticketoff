package repositories

import (
	"context"
	"fmt"
	"ticketoff/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id primitive.ObjectID) error
	GetUserByEmail(email string) (*models.User, error)
	GetUsers() ([]models.User, error)
	ConfirmEmail(email string) error
	VerifyUser(email string) error
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) CreateUser(user *models.User) error {
	collection := repo.db.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

func (repo *userRepository) VerifyUser(email string) error {
	collection := repo.db.Collection("users")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"email": email},
		bson.M{"$set": bson.M{"email_confirmed": true}},
	)
	return err
}

func (u *userRepository) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	collection := u.db.Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return &user, err
}

func (u userRepository) UpdateUser(updatedUser *models.User) (*models.User, error) {
	collection := u.db.Collection("users")
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": updatedUser.ID}, bson.M{"$set": updatedUser})
	return updatedUser, err
}

func (u userRepository) DeleteUser(id primitive.ObjectID) error {
	collection := u.db.Collection("users")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (u userRepository) GetUserByEmail(email string) (*models.User, error) {
	collection := u.db.Collection("users")
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return &user, err
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

func (u userRepository) VerifyUserEmail(email string) error {
	collection := u.db.Collection("users")
	_, err := collection.UpdateOne(context.Background(), bson.M{"email": email}, bson.M{"$set": bson.M{"verified": true}})
	return err
}

func (repo *userRepository) ConfirmEmail(email string) error {
	collection := repo.db.Collection("users")
	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"email": email},
		bson.M{"$set": bson.M{"email_confirmed": true}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with email %s", email)
	}
	return nil
}
