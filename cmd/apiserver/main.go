package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"regexp"
	"ticketoff/migrations"
	"ticketoff/models"
)

var db *gorm.DB

func init() {
	var err error
	db, err = migrations.InitDB("user=asyl password=1234 dbname=ticketoffdb host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	models.Migrate(db)
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}

/*
// Helper function to hash passwords
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Helper function to check if password is valid
func validatePassword(password string) bool {
	return len(password) >= 8
}
*/
/* Handler for creating a new user (POST)
func createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate user input
	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if !validatePassword(user.Password) {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Hash password before saving
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Create user
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}*/

func createUser(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a map to check for unexpected keys
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check for required fields: "Email" and "Password"
	if _, ok := body["Email"]; !ok {
		http.Error(w, "'Email' field is required", http.StatusBadRequest)
		return
	}
	if _, ok := body["Password"]; !ok {
		http.Error(w, "'Password' field is required", http.StatusBadRequest)
		return
	}

	// Validate that no unexpected fields are present
	allowedKeys := map[string]bool{
		"Email":    true,
		"Password": true,
	}

	for key := range body {
		if !allowedKeys[key] {
			http.Error(w, "Unexpected field: "+key, http.StatusBadRequest)
			return
		}
	}

	// Convert to User struct
	var user models.User
	user.Email = body["Email"].(string)
	user.Password = body["Password"].(string)

	// Allow empty Email and Password but check if they are both empty
	if user.Email == "" && user.Password == "" {
		// No error, as both fields can be empty
	} else {
		// Validate email if provided
		if user.Email != "" && !isValidEmail(user.Email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		// Validate password length if provided
		if user.Password != "" && len(user.Password) < 8 {
			http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}
	}

	// Save user to the database
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Handler for updating a user (PUT)
func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}

// Handler for deleting a user (DELETE)
func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	db.Delete(&user)
	w.WriteHeader(http.StatusNoContent)
}

// Handler for fetching all users (GET)
func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Handler for getting a user by ID (GET)
func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handleRequest() {
	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequest()
}
