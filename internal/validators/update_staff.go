package validators

import (
	"encoding/json"
	"time"

	"github.com/Confialink/wallet-pkg-utils/pointer"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UpdateStaffValidator struct {
	Data struct {
		Email       string `json:"email" binding:"required,email,uniqueEmail"`
		FirstName   string `json:"firstName" binding:"required,max=255"`
		LastName    string `json:"lastName" binding:"required,max=255"`
		MiddleName  string `json:"middleName" binding:"max=255"`
		PhoneNumber string `json:"phoneNumber" binding:"required,phonenumber,uniquePhoneNumber"`
		// User Details
		ClassId                    json.Number `json:"classId" binding:"required"`
		CountryOfResidenceIsoTwo   string      `json:"countryOfResidenceIsoTwo" binding:"max=2"`
		CountryOfCitizenshipIsoTwo string      `json:"countryOfCitizenshipIsoTwo" binding:"max=2"`
		DateOfBirth                *string     `json:"dateOfBirth"`
		// Addresses
		PhysicalAdressValidator
	} `json:"data"`
	UserModel models.User `json:"-"`
}

func GetUpdateStaffValidator() UpdateStaffValidator {
	return UpdateStaffValidator{}
}

// GetUpdateUserValidatorFillWith ...
func GetUpdateStaffValidatorFillWith(user models.User) UpdateUserValidator {

	v := GetUpdateUserValidator()
	v.Data.UID = user.UID
	v.Data.Username = user.Username
	v.Data.Email = user.Email
	v.Data.FirstName = user.FirstName
	v.Data.LastName = user.LastName
	v.Data.MiddleName = user.MiddleName
	v.Data.PhoneNumber = user.PhoneNumber
	v.Data.SmsPhoneNumber = user.SmsPhoneNumber
	v.Data.IsCorporate = *user.IsCorporate
	v.Data.RoleName = user.RoleName
	v.Data.Status = user.Status
	v.Data.Position = user.Position
	v.Data.InternalNotes = user.InternalNotes
	v.Data.ClassId = user.ClassId

	v.Data.CompanyName = user.CompanyDetails.CompanyName
	v.Data.CompanyType = user.CompanyDetails.CompanyType
	v.Data.CompanyRole = user.CompanyDetails.CompanyRole
	v.Data.DirectorFirstName = user.CompanyDetails.DirectorFirstName
	v.Data.DirectorLastName = user.CompanyDetails.DirectorLastName

	//v.User.DateOfBirth = user.DateOfBirth // TODO: add validation
	v.Data.DocumentType = user.GetDocumentType()
	v.Data.DocumentPersonalId = user.DocumentPersonalId
	v.Data.CountryOfResidenceIsoTwo = user.CountryOfResidenceIsoTwo
	v.Data.CountryOfCitizenshipIsoTwo = user.CountryOfCitizenshipIsoTwo

	v.Data.HomePhoneNumber = user.HomePhoneNumber
	v.Data.OfficePhoneNumber = user.OfficePhoneNumber
	v.Data.Fax = user.Fax

	// Copy user group id pointer
	if user.UserGroupId != nil {
		v.Data.UserGroupId = pointer.ToUint64(*user.UserGroupId)
	}

	v.UserModel = user

	return v
}

func (s *UpdateStaffValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.UserModel.Email = s.Data.Email
	s.UserModel.FirstName = s.Data.FirstName
	s.UserModel.LastName = s.Data.LastName
	s.UserModel.MiddleName = s.Data.MiddleName

	s.UserModel.PhoneNumber = s.Data.PhoneNumber
	s.UserModel.ClassId = s.Data.ClassId
	s.UserModel.LastActedAt = time.Now()

	//s.UserModel.DateOfBirth = s.User.DateOfBirth // TODO: add validation

	return nil
}
