package routes

import (
	db "knockNSell/db/gen"
	logger "knockNSell/logger"
	"net/http"
	"strings"
	"time"

	helper "knockNSell/helpers"

	"github.com/gin-gonic/gin"
)

func (s *Server) LoginUser(c *gin.Context) {
	var payload PhoneModel

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.Request = c.Request.WithContext(
			logger.SetLogMessageAndFields(c.Request.Context(), "🚨 Could not map ", gin.H{
				"error": error.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}

	dbResponse, error := s.q.GetUserByPhoneNumber(c.Request.Context(), payload.PhoneNumber)
	if error != nil {
		c.Request = c.Request.WithContext(
			logger.SetLogMessageAndFields(c.Request.Context(), "🚨 User not found for phone number :- "+payload.PhoneNumber, gin.H{
				"error": error.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status_code": 401,
			"db message":  error.Error(),
			"message":     "We could not find your phone.",
		})
		return
	} else {
		var authTokenExpiresAt = time.Now().Add(24 * time.Hour)         // Access token expires in 24 hours
		var refreshTokenExpiresAt = time.Now().Add(14 * 24 * time.Hour) // Refresh token expires in 14 days

		authToken, refreshToken := helper.CreateAuthAndRefreshToken(authTokenExpiresAt, refreshTokenExpiresAt, dbResponse)

		payLoad := db.CreateAuthTokenParams{
			UserID:                dbResponse.ID,
			AuthToken:             authToken,
			RefreshToken:          refreshToken,
			UserAgent:             helper.ToNullString(c.GetHeader("User-Agent")),
			AuthTokenExpiresAt:    helper.ToNullTime(authTokenExpiresAt),
			RefreshTokenExpiresAt: helper.ToNullTime(refreshTokenExpiresAt),
			IpAddress:             helper.ToNullString(c.Request.RemoteAddr),
		}

		dbAuth, error := s.q.CreateAuthToken(c.Request.Context(), payLoad)
		if error != nil {
			c.Request = c.Request.WithContext(
				logger.SetLogMessageAndFields(c.Request.Context(), "🚨 Could not save the tokens to Database ", gin.H{
					"error": error.Error(),
				}),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status_code": 401,
				"message":     "Could not save the tokens to Database",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status_code":   200,
				"message":       "Logged In Successfully",
				"auth_token":    dbAuth.AuthToken,
				"refresh_token": dbAuth.RefreshToken,
			})
		}
	}
}

func (s *Server) SignUpUser(c *gin.Context) {

	type userSingUpModel struct {
		PhoneNumber      string `json:"phoneNumber" binding:"required"`
		AccountType      string `json:"accountType"`
		Email            string `json:"email,omitempty"`
		Name             string `json:"name,omitempty"`
		Photo            string `json:"photo,omitempty"`
		Gender           string `json:"gender,omitempty"`
		AadharNumber     string `json:"aadharNumber,omitempty"`
		AadharPhotoFront string `json:"aadharPhotoFront,omitempty"`
		AadharPhotoBack  string `json:"aadharPhotoBack,omitempty"`
		VehicleType      string `json:"vehicleType,omitempty"`
		Age              int    `json:"age,omitempty"`
		GstNumber        string `json:"gstNumber,omitempty"`
	}

	var payload userSingUpModel

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.Request = c.Request.WithContext(
			logger.SetLogMessageAndFields(c.Request.Context(), "🚨 Could not save the tokens to Database ", gin.H{
				"error": error.Error(),
			}),
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}
	dbResponse, error := s.q.CreateUser(c.Request.Context(), db.CreateUserParams{
		PhoneNumber: payload.PhoneNumber,
		AccountType: payload.AccountType,
		Email:       helper.ToNullString(payload.Email),
	})

	if error != nil {
		if strings.Contains(error.Error(), "users_phone_number_key") {
			c.Request = c.Request.WithContext(
				logger.SetLogMessageAndFields(c.Request.Context(), "🚨 The user already exits with phone :- "+payload.PhoneNumber, gin.H{
					"error": error.Error(),
				}),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status_code": 401,
				"db message":  error.Error(),
				"message":     "The user already exist.",
			})
		} else {
			c.Request = c.Request.WithContext(
				logger.SetLogMessageAndFields(c.Request.Context(), "🚨 Can not create a user with phone number :- "+payload.PhoneNumber, gin.H{
					"error": error.Error(),
				}),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status_code": 400,
				"db message":  error.Error(),
				"message":     "Can not save the number.",
			})
		}
		return
	} else {

		var authTokenExpiresAt = time.Now().Add(24 * time.Hour)         // Access token expires in 24 hours
		var refreshTokenExpiresAt = time.Now().Add(14 * 24 * time.Hour) // Refresh token expires in 14 days

		authToken, refreshToken := helper.CreateAuthAndRefreshToken(authTokenExpiresAt, refreshTokenExpiresAt, dbResponse)

		payLoad := db.CreateAuthTokenParams{
			UserID:                dbResponse.ID,
			AuthToken:             authToken,
			RefreshToken:          refreshToken,
			UserAgent:             helper.ToNullString(c.GetHeader("User-Agent")),
			AuthTokenExpiresAt:    helper.ToNullTime(authTokenExpiresAt),
			RefreshTokenExpiresAt: helper.ToNullTime(refreshTokenExpiresAt),
			IpAddress:             helper.ToNullString(c.Request.RemoteAddr),
		}

		dbAuth, error := s.q.CreateAuthToken(c.Request.Context(), payLoad)
		if error != nil {
			c.Request = c.Request.WithContext(
				logger.SetLogMessageAndFields(c.Request.Context(), "🚨 Could not save the tokens to Database.", gin.H{
					"error": error.Error(),
				}),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status_code": 401,
				"message":     "Could not save the tokens to Database",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status_code":   200,
				"message":       "Account created successfully",
				"auth_token":    dbAuth.AuthToken,
				"refresh_token": dbAuth.RefreshToken,
			})
		}
	}
}
