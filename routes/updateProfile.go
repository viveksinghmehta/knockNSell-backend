package routes

import (
	"fmt"
	db "knockNSell/db/gen"
	helper "knockNSell/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) UpdateProfile(c *gin.Context) {

	type userSingUpModel struct {
		PhoneNumber      string `json:"phoneNumber" binding:"required"`
		AccountType      string `json:"accountType,omitempty"`
		Email            string `json:"email,omitempty"`
		Name             string `json:"name,omitempty"`
		Photo            string `json:"photo,omitempty"`
		Gender           string `json:"gender,omitempty"`
		AadharNumber     string `json:"aadharNumber,omitempty"`
		AadharPhotoFront string `json:"aadharPhotoFront,omitempty"`
		AadharPhotoBack  string `json:"aadharPhotoBack,omitempty"`
		VehicleType      string `json:"vehicleType,omitempty"`
		Age              int32  `json:"age,omitempty"`
		GstNumber        string `json:"gstNumber,omitempty"`
	}

	var payload userSingUpModel

	if error := c.ShouldBindJSON(&payload); error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}
	fmt.Printf("Updating user with phoneNumber=%s accountType=%s email=%s\n",
		payload.PhoneNumber,
		payload.AccountType,
		payload.Email,
	)

	_, error := s.q.UpdateUserByPhoneNumber(c.Request.Context(), db.UpdateUserByPhoneNumberParams{
		PhoneNumber: payload.PhoneNumber,
		AccountType: payload.AccountType,
		Email:       helper.ToNullString(payload.Email),
		Name:        payload.Name,
		Photo:       helper.ToNullString(payload.Photo),
		Gender: db.NullGenderEnum{
			GenderEnum: db.GenderEnumMale,
		},
		AadharNumber:     helper.ToNullString(payload.AadharNumber),
		AadharPhotoFront: helper.ToNullString(payload.AadharPhotoFront),
		AadharPhotoBack:  helper.ToNullString(payload.AadharPhotoBack),
		VehicleType:      helper.ToNullString(payload.VehicleType),
		Age:              helper.ToNullInt(payload.Age),
		GstNumber:        helper.ToNullString(payload.GstNumber),
	})

	if error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status_code": 401,
			"db message":  error.Error(),
			"message":     "Can not save the details, please try again.",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusOK,
			"message":     "Profile updated.",
		})
	}
}
