package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
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

// Handler for creating a new user (POST)
func createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
	if err := db.First(&models.User{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	db.Delete(&user)
	w.WriteHeader(http.StatusNoContent)
}

// Handler for the home page (GET)
func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the TicketOff API!"))
}

func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Define home page route
	router.HandleFunc("/", homePage).Methods("GET")

	// Start the server with the router handling requests
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequest()
}
