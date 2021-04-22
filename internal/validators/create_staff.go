package validators

import (
	"encoding/json"
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CreateStaffValidator struct {
	Data struct {
		// Username        string `json:"username" binding:"omitempty,min=4,max=50,usernameChars"`
		Email           string `json:"email" binding:"required,email,uniqueEmail"`
		Password        string `json:"password" binding:"required,min=8,max=255,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
		FirstName       string `json:"firstName" binding:"required,max=255"`
		LastName        string `json:"lastName" binding:"required,max=255"`
		MiddleName      string `json:"middleName" binding:"max=255"`
		PhoneNumber     string `json:"phoneNumber" binding:"required,phonenumber,uniquePhoneNumber"`
		RoleName        string `json:"roleName" binding:"required,roleoneof"`
		// User Details
		ClassId                  json.Number `json:"classId" binding:"required"`
		CountryOfResidenceIso2   string      `json:"countryOfResidenceISO2" binding:"max=2"`
		CountryOfCitizenshipIso2 string      `json:"countryOfCitizenshipISO2" binding:"max=2"`
		DateOfBirthYear          uint64      `json:"dateOfBirthYear"`
		DateOfBirthMonth         uint64      `json:"dateOfBirthMonth"`
		DateOfBirthDay           uint64      `json:"dateOfBirthDay"`
		UserGroupId              *uint64     `json:"userGroupId,omitempty"`
		CompanyDetails
		// Addresses
		PhysicalAdressValidator
	} `json:"data"`
	UserModel models.User `json:"-"`
}

// BindJSON binding from JSON
func (s *CreateStaffValidator) BindJSON(ctx *gin.Context) error {
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
	s.UserModel.RoleName = s.Data.RoleName
	s.UserModel.UserGroupId = s.Data.UserGroupId

	s.UserModel.ClassId = s.Data.ClassId
	s.UserModel.LastActedAt = time.Now()

	// Company details
	s.UserModel.CompanyDetails.CompanyName = s.Data.CompanyName
	s.UserModel.CompanyDetails.CompanyType = s.Data.CompanyType
	s.UserModel.CompanyDetails.CompanyRole = s.Data.CompanyRole
	s.UserModel.CompanyDetails.DirectorFirstName = s.Data.DirectorFirstName
	s.UserModel.CompanyDetails.DirectorLastName = s.Data.DirectorLastName

	return nil
}
