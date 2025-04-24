package routes

import "github.com/gin-gonic/gin"

type PhoneModel struct {
	phoneNumber string
}

func Sendotp(c *gin.Context) {
	var requestBody PhoneModel
	if error := c.BindJSON(&requestBody); error != nil {
		c.JSON(401, gin.H{
			"message": "There was some error",
		})
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
