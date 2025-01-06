package utils

import "regexp"

// IsValidEmail checks email format
func IsValidEmail(email string) bool {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}

// ValidatePassword ensures password is strong
func ValidatePassword(password string) bool {
	return len(password) >= 8
}
