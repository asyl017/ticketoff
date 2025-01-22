package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
	_ "gopkg.in/gomail.v2"
	"log"
	"net/http"
	"ticketoff/handler"
	"ticketoff/migrations"
	"ticketoff/models"
	"ticketoff/repositories"
	"ticketoff/utils"
)

var Dialer *gomail.Dialer

func main() {
	utils.InitLogger()
	router := mux.NewRouter()

	router.Use(utils.LoggingMiddleware)

	db := InitDB()

	db.LogMode(true)

	userRepo := repositories.NewUserRepository(db)
	filmRepo := repositories.NewFilmRepository(db)
	userRouter := handler.NewUserRouter(userRepo)
	authHandler := handler.NewAuthHandler(userRepo)
	filmHandler := handler.NewFilmHandler(filmRepo)
	// Correctly reference handler
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

func InitDB() *gorm.DB {
	db, err := migrations.InitDB("user=asyl password=1234 dbname=ticketoffdb host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	models.MigrateUser(db)
	if err := db.AutoMigrate(&models.Film{}).Error; err != nil {
		log.Fatalf("Film migration failed: %v", err)
	}
	return db
}
