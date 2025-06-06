package helper

import (
	model "knockNSell/db/gen"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
func GenerateAccessToken(user model.User, expiresAt time.Time) (string, error) {
	// Set custom and registered claims
	claims := CustomClaims{
		UserID:      user.ID,
		Name:        user.Name,
		Mobile:      user.PhoneNumber,
		AccountType: user.AccountType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
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
func GenerateRefreshToken(user model.User, expiresAt time.Time) (string, error) {
	// Set custom and registered claims
	claims := CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
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

func CreateAuthAndRefreshToken(authExpiresAt time.Time, refreshExpiresAt time.Time, user model.User) (authToken, refreshToken string) {
	authToken, authError := GenerateAccessToken(user, authExpiresAt)
	refreshToken, refreshError := GenerateRefreshToken(user, refreshExpiresAt)
	if authError != nil && refreshError != nil {
		log.WithFields(log.Fields{
			"authErrorMessage":    authError.Error(),
			"refreshErrorMessage": refreshError.Error(),
		}).Error("Could not generate the token." + "🚨")
	}
	return authToken, refreshToken
}
