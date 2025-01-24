package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"ticketoff/handler"
	"ticketoff/migrations"
	"ticketoff/repositories"
	"ticketoff/utils"
)

func main() {
	utils.InitLogger()
	router := mux.NewRouter()

	router.Use(utils.LoggingMiddleware)

	db := InitDB()

	userRepo := repositories.NewUserRepository(db)
	filmRepo := repositories.NewFilmRepository(db)
	userRouter := handler.NewUserRouter(userRepo)
	authHandler := handler.NewAuthHandler(userRepo)
	filmHandler := handler.NewFilmHandler(filmRepo)

	router.HandleFunc("/users", userRouter.CreateUser).Methods("POST")
	router.HandleFunc("/users", userRouter.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userRouter.GetUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", userRouter.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userRouter.DeleteUser).Methods("DELETE")

	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	router.HandleFunc("/films", filmHandler.GetFilms).Methods("GET")
	router.HandleFunc("/films/{id}", filmHandler.GetFilmByID).Methods("GET")

	router.HandleFunc("/send-email", handler.SendEmail).Methods("POST")

	emailHandler := handler.EmailHandler{UserRepo: userRepo}
	router.HandleFunc("/confirm-email", emailHandler.ConfirmEmail).Methods("GET")

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)
	log.Fatal(http.ListenAndServe(":8080", cors(router)))
}

func InitDB() *mongo.Database {
	db, err := migrations.InitDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	return db
}
