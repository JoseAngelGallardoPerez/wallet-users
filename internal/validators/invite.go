package validators

import (
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type InviteValidator struct {
	Data struct {
		Email            string `json:"email" binding:"required,email,uniqueEmail"`
		FirstName        string `json:"firstName" binding:"required,max=255"`
		LastName         string `json:"lastName" binding:"required,max=255"`
		MiddleName       string `json:"middleName" binding:"max=255"`
		PhoneNumber      string `json:"phoneNumber" binding:"required,phonenumber,uniquePhoneNumber"`
		DateOfBirthYear  uint64 `json:"dateOfBirthYear"`
		DateOfBirthMonth uint64 `json:"dateOfBirthMonth"`
		DateOfBirthDay   uint64 `json:"dateOfBirthDay"`
		Password         string `json:"password" binding:"required,min=8,max=255,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`
		ConfirmPassword  string `json:"confirmPassword" binding:"required,eqfield=Password"`
		// Addresses
		PhysicalAdressValidator
	} `json:"data"`
	UserModel models.User `json:"-"`
}

// BindJSON binding from JSON
func (s *InviteValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	// s.UserModel.Username = s.User.Username
	s.UserModel.Email = s.Data.Email
	s.UserModel.FirstName = s.Data.FirstName
	s.UserModel.LastName = s.Data.LastName
	s.UserModel.MiddleName = s.Data.MiddleName

	s.UserModel.Password = s.Data.Password
	s.UserModel.PhoneNumber = s.Data.PhoneNumber
	s.UserModel.LastActedAt = time.Now()

	//s.UserModel.DateOfBirthYear = s.User.DateOfBirthYear
	//s.UserModel.DateOfBirthDay = s.User.DateOfBirthDay
	//s.UserModel.DateOfBirthMonth = s.User.DateOfBirthMonth

	return nil
}
