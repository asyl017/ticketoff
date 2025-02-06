package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"ticketoff/repositories"
	"ticketoff/utils"
)

func AdminMiddleware(userRepo repositories.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("AdminMiddleware: Checking for token cookie")
			c, err := r.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					log.Println("AdminMiddleware: No token cookie found, redirecting to login")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				log.Println("AdminMiddleware: Error retrieving token cookie:", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			tokenStr := c.Value
			log.Println("AdminMiddleware: Parsing token")
			email, err := utils.ParseToken(tokenStr)
			if err != nil {
				log.Println("AdminMiddleware: Error parsing token:", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			log.Println("AdminMiddleware: Retrieving user by email:", email)
			user, err := userRepo.GetUserByEmail(email)
			if err != nil {
				log.Println("AdminMiddleware: Error retrieving user:", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !user.IsAdmin {
				log.Println("AdminMiddleware: User is not an admin, redirecting to login")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			log.Println("AdminMiddleware: User is an admin, proceeding with request")
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
