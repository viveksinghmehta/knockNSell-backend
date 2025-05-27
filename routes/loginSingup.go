package routes

import (
	"database/sql"
	db "knockNSell/db/gen"
	"net/http"
	"strings"
	"time"

	helper "knockNSell/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Server) LoginUser(c *gin.Context) {
	var payload PhoneModel
	start := time.Now()

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		log.WithFields(helper.GetExtraFieldsForSlackLog(c, start)).Error(error.Error())
		return
	}

	dbResponse, error := s.q.GetUserByPhoneNumber(c.Request.Context(), payload.PhoneNumber)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"db message": error.Error(),
			"message":    "We could not find your phone.",
		})
		log.WithFields(helper.GetExtraFieldsForSlackLog(c, start)).Error(error.Error() + "🚨")
		return
	} else {
		var authTokenExpiresAt = time.Now().Add(24 * time.Hour)         // Access token expires in 24 hours
		var refreshTokenExpiresAt = time.Now().Add(14 * 24 * time.Hour) // Refresh token expires in 14 days
		authToken, authError := helper.GenerateAccessToken(dbResponse, authTokenExpiresAt)
		refreshoken, refreshError := helper.GenerateRefreshToken(dbResponse, refreshTokenExpiresAt)
		if authError != nil && refreshError != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status_code": 401,
				"message":     "Could not generate the token.",
			})
			log.WithFields(helper.GetExtraFieldsForSlackLog(c, start)).Error("Could not generate the token." + "🚨")
			return
		} else {
			payLoad := db.CreateAuthTokenParams{
				UserID:       dbResponse.ID,
				AuthToken:    authToken,
				RefreshToken: refreshoken,
				UserAgent: sql.NullString{
					String: c.GetHeader("User-Agent"),
					Valid:  true,
				},
				AuthTokenExpiresAt: sql.NullTime{
					Time:  authTokenExpiresAt,
					Valid: true,
				},
				RefreshTokenExpiresAt: sql.NullTime{
					Time:  refreshTokenExpiresAt,
					Valid: true,
				},
				IpAddress: sql.NullString{
					String: c.Request.RemoteAddr,
					Valid:  true,
				},
			}
			dbAuth, error := s.q.CreateAuthToken(c.Request.Context(), payLoad)
			if error != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status_code": 401,
					"message":     "Could not save the tokens to Database",
				})
				log.WithFields(helper.GetExtraFieldsForSlackLog(c, start)).Error("Could not save the tokens to Database." + "🚨")
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

}

func (s *Server) SignUpUser(c *gin.Context) {
	type userSingUpModel struct {
		PhoneNumber string `json:"phoneNumber" binding:"required"`
		AccountType string `json:"accountType"`
	}

	var payload userSingUpModel

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}

	dbResponse, error := s.q.CreateUser(c.Request.Context(), db.CreateUserParams{
		PhoneNumber: payload.PhoneNumber,
		AccountType: payload.AccountType,
	})

	if error != nil {
		if strings.Contains(error.Error(), "users_phone_number_key") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":     401,
				"db message": error.Error(),
				"message":    "The user already exist.",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":     400,
				"db message": error.Error(),
				"message":    "Can not save the number.",
			})
		}
		return
	} else {

		type UserModel struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			AccountType  string `json:"accountType"`
			PhoneNumber  string `json:"phoneNumber"`
			Email        string `json:"email"`
			Photo        string `json:"photo"`
			AadharNumber string `json:"aadharNumber"`
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"phone":   dbResponse.PhoneNumber,
			"message": "Logged In Successfully",
			"user": UserModel{
				ID:           dbResponse.ID.String(),
				Name:         dbResponse.Name,
				AccountType:  dbResponse.AccountType,
				PhoneNumber:  dbResponse.PhoneNumber,
				Email:        dbResponse.Email.String,
				Photo:        dbResponse.Photo.String,
				AadharNumber: dbResponse.AadharNumber.String,
			},
			"db message": dbResponse.ID,
		})
	}
}
