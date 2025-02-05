package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"ticketoff/models"
	"ticketoff/repositories"
	"ticketoff/utils"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "random"
)

type RegistrationHandler struct {
	UserRepo repositories.UserRepository
}

func (h *RegistrationHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.EmailConfirmed = false

	err = h.UserRepo.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	verificationLink := fmt.Sprintf("http://localhost:8080/verify?email=%s", user.Email)
	subject := "Email Verification"
	body := "Please verify your email by clicking on the following link: " + verificationLink
	err = utils.SendVerificationEmail(user.Email, subject, body)
	if err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Registration successful, please check your email to verify your account"))
}

func NewRegistrationHandler(userRepo repositories.UserRepository) *RegistrationHandler {
	return &RegistrationHandler{UserRepo: userRepo}
}

func (h *RegistrationHandler) HandleGoogleRegister(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *RegistrationHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(oauth2.NoContext, token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer userInfo.Body.Close()

	var user models.User
	if err := json.NewDecoder(userInfo.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle user registration or login with the obtained user info
	_, err = h.UserRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User does not exist, create a new user
			newUser := models.User{
				Email:    user.Email,
				Password: "", // No password for OAuth users
			}
			if err := h.UserRepo.CreateUser(&newUser); err != nil {
				http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Failed to retrieve user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Generate JWT token for the user
	jwtToken, err := utils.GenerateJWT(&user)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtToken,
		Expires: time.Now().Add(72 * time.Hour),
	})

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func (h *RegistrationHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Extract the verification token from the request
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Verification token is missing", http.StatusBadRequest)
		return
	}

	// Verify the token (this is a placeholder, implement your own verification logic)
	/*email, err := utils.VerifyEmailToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}*/

	// Update the user's email verification status in the database
	/*err = h.UserRepo.VerifyUserEmail(email)
	if err != nil {
		http.Error(w, "Failed to verify email: "+err.Error(), http.StatusInternalServerError)
		return
	}*/

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully"))
}

func (h *RegistrationHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	_, err := h.UserRepo.GetUserByEmail(user.Email)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash the password
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(password)
	user.EmailConfirmed = false

	// Create a new user
	if err := h.UserRepo.CreateUser(&user); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a verification token
	verificationToken := utils.GenerateToken(user.Email)

	// Send verification email
	verificationLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", verificationToken)
	subject := "Email Verification"
	body := "Please verify your email by clicking on the following link: " + verificationLink
	err = utils.SendVerificationEmail(user.Email, subject, body)
	if err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Registration successful, please check your email to verify your account"))
}
