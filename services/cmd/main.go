package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"ticketoff/handler"
	"ticketoff/migrations"
	"ticketoff/models"
	"ticketoff/repositories"
)

func main() {
	router := mux.NewRouter()

	db := InitDB()
	userRepo := repositories.NewUserRepository(db)
	userRouter := handler.NewUserRouter(userRepo)
	authHandler := handler.NewAuthHandler(userRepo)
	// Correctly reference handler
	router.HandleFunc("/users", userRouter.CreateUser).Methods("POST")
	router.HandleFunc("/users", userRouter.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userRouter.GetUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", userRouter.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userRouter.DeleteUser).Methods("DELETE")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)
	log.Fatal(http.ListenAndServe(":8080", cors(router)))
}

func InitDB() *gorm.DB {
	db, err := migrations.InitDB("user=asyl password=1234 dbname=ticketoffdb host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	models.Migrate(db)
	return db
}
