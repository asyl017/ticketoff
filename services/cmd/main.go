package main

import (
	"fmt"
	"log"
	"net/http"
	"ticketoff/handler"
	"ticketoff/middleware"
	"ticketoff/repositories"
	"ticketoff/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

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
	registrationHandler := handler.NewRegistrationHandler(userRepo)
	emailHandler := handler.EmailHandler{UserRepo: userRepo}

	router.HandleFunc("/users", userRouter.CreateUser).Methods("POST")
	router.HandleFunc("/users", userRouter.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userRouter.GetUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", userRouter.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userRouter.DeleteUser).Methods("DELETE")

	router.HandleFunc("/login",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Login endpoint called")
			http.ServeFile(w, r, "templates/login.html")
		}).Methods("GET")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	router.HandleFunc("/films", filmHandler.GetFilms).Methods("GET")
	router.HandleFunc("/films/{id}", filmHandler.GetFilmByID).Methods("GET")

	router.HandleFunc("/send-email", handler.SendEmail).Methods("POST")

	router.HandleFunc("/confirm-email", emailHandler.ConfirmEmail).Methods("GET")

	// Serve the registration HTML page
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// OAuth 2.0 registration routes
	router.HandleFunc("/register", registrationHandler.Register).Methods("POST")
	router.HandleFunc("/callback", registrationHandler.HandleGoogleCallback)

	// Email verification routes
	router.HandleFunc("/register", registrationHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Verify endpoint called")
		handler.VerifyUser(w, r, db)
	}).Methods("GET")
	// Protected routes
	protected := router.PathPrefix("/admin").Subrouter()
	protected.HandleFunc("/dashboard", adminDashboard).Methods("GET")

	// Protected routes
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AdminMiddleware(userRepo))
	adminRouter.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/admin.html")
	}).Methods("GET")
	adminRouter.HandleFunc("/admin_panel", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/admin_panel.html")
	}).Methods("GET")
	adminRouter.HandleFunc("/dashboard", adminDashboard).Methods("GET")

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)
	log.Fatal(http.ListenAndServe(":8080", cors(router)))

}

func InitDB() *mongo.Database {
	// MongoDB Atlas URI (Replace this with your actual MongoDB Atlas URI)
	const uri = "mongodb+srv://asyl17:1234567654321@cluster0.t7gd8.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	// Set MongoDB Atlas connection options
	clientOptions := options.Client().ApplyURI(uri)

	// Create a MongoDB client
	client, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Fatal("Error initializing MongoDB client: ", err)
	}

	// Ping the MongoDB server to check if the connection is successful
	err = client.Ping(nil, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB: ", err)
	}

	log.Println("Database connected successfully to MongoDB Atlas")

	// Return the database instance
	db := client.Database("ticketoffdb") // Use your database name here
	return db
}
func adminDashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the admin dashboard!")
}
