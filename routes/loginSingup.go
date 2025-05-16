package routes

import (
	db "knockNSell/db/gen"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) LoginUser(c *gin.Context) {
	var payload PhoneModel

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}

	dbResponse, error := s.q.GetUserByPhoneNumber(c.Request.Context(), payload.PhoneNumber)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"db message": error.Error(),
			"message":    "We could not find your phone.",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"phone":      dbResponse.PhoneNumber,
			"message":    "Logged In Successfully",
			"user":       dbResponse,
			"db message": dbResponse.ID,
		})
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

	// // Check for User in DB from phone number
	// _, error := s.q.GetUserByPhoneNumber(c.Request.Context(), payload.PhoneNumber)
	// if error == nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"db message": error.Error(),
	// 		"message":    "User already exist",
	// 	})
	// 	return
	// } else {
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
	// }
}
