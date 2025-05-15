package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type VerifyOTPModel struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

func containsString(arr pq.StringArray, str string) bool {
	for _, element := range arr {
		if element == str {
			return true
		}
	}
	return false
}

func (s *Server) VerifyOTP(c *gin.Context) {
	var payload VerifyOTPModel

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}

	dbResponse, error := s.q.GetOTPByPhoneNumber(c.Request.Context(), payload.PhoneNumber)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"db message": error.Error(),
			"message":    "Can't get value of " + payload.PhoneNumber + " from the server",
		})
		return
	} else {
		otps := dbResponse.Otp
		if containsString(otps, payload.OTP) {
			c.JSON(http.StatusAccepted, gin.H{
				"success":     true,
				"message":     "OTP successfully verified.",
				"id":          dbResponse.ID,
				"phoneNumber": dbResponse.PhoneNumber,
			})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Could not verify OTP please try again",
			})
		}
	}
}
