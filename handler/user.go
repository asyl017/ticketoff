package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"ticketoff/models"
	"ticketoff/repositories"
	"ticketoff/utils"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UserRouter interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type userRouter struct {
	userRepo repositories.UserRepository
}

func NewUserRouter(userRepo repositories.UserRepository) UserRouter {
	return &userRouter{
		userRepo: userRepo,
	}
}

// CreateUser (Sign-Up Handler)
func (u userRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	// Validate user input
	if user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Email and password are required"})
		return
	}
	if !utils.IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Invalid email format"})
		return
	}
	if !utils.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Password must be at least 8 characters"})
		return
	}

	// Check if user already exists
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusConflict, Message: "User already exists"})
		return
	}

	// Generate a new ObjectID for the user
	user.ID = primitive.NewObjectID()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusInternalServerError, Message: "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	// Save user
	err = u.userRepo.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusConflict, Message: "Error creating user: " + err.Error()})
		return
	}

	// Send confirmation email
	confirmationLink := fmt.Sprintf("http://localhost:8080/confirm-email?token=%s", utils.GenerateToken(user.Email))
	err = utils.SendEmail(user.Email, "Confirm your email", "Please confirm your email by clicking the following link: "+confirmationLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusInternalServerError, Message: "Error sending confirmation email: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUsers
func (u userRouter) GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all users")

	users, err := u.userRepo.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusInternalServerError, Message: "Error fetching users"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID
func (u userRouter) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Invalid user ID"})
		return
	}

	user, err := u.userRepo.GetUserByID(objectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(HttpError{Code: http.StatusNotFound, Message: "User not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusInternalServerError, Message: "Internal server error"})
		log.Println("Error fetching user by ID:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser
func (u userRouter) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Invalid user ID"})
		return
	}

	var updatedUser models.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusInternalServerError, Message: "Error hashing password"})
		return
	}
	updatedUser.Password = string(hashedPassword)
	updatedUser.ID = objectID

	updatedUserPtr, err := u.userRepo.UpdateUser(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusInternalServerError, Message: "Error updating user: " + err.Error()})
		return
	}

	response := map[string]interface{}{
		"unhashed_password": updatedUser.Password,
		"user":              updatedUserPtr,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteUser
func (u userRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Invalid user ID"})
		return
	}

	err = u.userRepo.DeleteUser(objectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(HttpError{Code: http.StatusNotFound, Message: "User not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusInternalServerError, Message: "Error deleting user: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u userRouter) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusBadRequest, Message: "Email is required"})
		return
	}

	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(HttpError{Code: http.StatusNotFound, Message: "User not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
