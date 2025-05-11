package routes

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"

	db "knockNSell/db/gen"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	openApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type PhoneModel struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

func get6DigitNumber() string {
	max := big.NewInt(900000) // Upper bound (exclusive)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	code := n.Int64() + 100000 // Shift to range [100000, 999999]
	log.Printf("Secure 6-digit code: %d\n", code)
	return strconv.FormatInt(code, 10)
}

func twillioClient(phone, message string) (*openApi.ApiV2010Message, error) {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	password := os.Getenv("TWILIO_AUTH_TOKEN")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: sid,
		Password: password,
	})

	params := &openApi.CreateMessageParams{}
	countryCode := "+91"
	params.SetFrom("+19786984267")    // Your Twilio phone number
	params.SetTo(countryCode + phone) // Recipient's phone number
	params.SetBody(message)

	return client.Api.CreateMessage(params)
}

func (s *Server) Sendotp(c *gin.Context) {
	var payload PhoneModel

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}

	code := get6DigitNumber()
	messageBody := "Your Knock & Sell 6 digit code is: " + code + ". Please do not share it."

	resp, error := twillioClient(payload.PhoneNumber, messageBody)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	} else {
		if resp.Sid != nil {
			log.Println(*resp.Sid)
		} else {
			log.Println(*resp.Sid)
		}
		dbResponse, error := s.q.UpsertOtpVerification(c.Request.Context(), db.UpsertOtpVerificationParams{
			PhoneNumber: payload.PhoneNumber,
			Otp:         code,
		})
		if error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"db message": error.Error(),
				"message":    "Can't enter OTP to the database",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"status":       resp.Status,
				"code":         code,
				"date_created": resp.DateCreated,
				"phone":        resp.To,
				"date_updated": resp.DateUpdated,
				"message":      "OTP send successfully",
				"db message":   dbResponse.ID,
			})
		}
	}

}
