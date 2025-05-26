package helper

import (
	model "knockNSell/db/gen"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims represents the JWT claims, embedding the standard claims
type CustomClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Mobile      string    `json:"mobile"`
	AccountType string    `json:"account_type"`
	jwt.RegisteredClaims
}

const appName = "knockNSell"

var accessSecret = []byte(os.Getenv("ACCESS_SECRET_KEY"))
var refreshSecret = []byte(os.Getenv("REFRESH_SECRET_KEY"))

// GenerateAccessToken creates a JWT access token with a short expiration time
func GenerateAccessToken(user model.User) (string, error) {
	// Set custom and registered claims
	claims := CustomClaims{
		UserID:      user.ID,
		Name:        user.Name,
		Mobile:      user.PhoneNumber,
		AccountType: user.AccountType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Access token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    appName,
			Subject:   user.ID.String(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString(accessSecret)
}

// GenerateRefreshToken creates a JWT refresh token with a longer expiration time
func GenerateRefreshToken(user model.User) (string, error) {
	// Set custom and registered claims
	claims := CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * 24 * time.Hour)), // Refresh token expires in 14 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    appName,
			Subject:   user.ID.String(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString(refreshSecret)
}
