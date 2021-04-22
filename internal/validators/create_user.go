package validators

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/Confialink/wallet-users/internal/db/models"
)

// CreateUserValidator containing two parts:
// - Validator: write the form/json checking rule according to the doc https://github.com/go-playground/validator
// - DataModel: fill with data from Validator after invoking ctx.ShouldBindWith(self)
// Then, you can just call CreateUserValidator.userModel after the data is ready in DataModel.
type CreateUserValidator struct {
	Data struct {
		Username        string  `json:"username" binding:"omitempty,min=4,max=50,usernameChars"`
		Email           string  `json:"email" binding:"required,email,uniqueEmail"`
		Password        string  `json:"password" binding:"required,min=8,max=255,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`
		ConfirmPassword string  `json:"confirmPassword" binding:"required,eqfield=Password"`
		FirstName       string  `json:"firstName" binding:"required,max=255"`
		LastName        string  `json:"lastName" binding:"required,max=255"`
		MiddleName      string  `json:"middleName" binding:"max=255"`
		Nickname        string  `json:"nickname" binding:"max=100"`
		PhoneNumber     string  `json:"phoneNumber" binding:"required,uniquePhoneNumber"`
		SmsPhoneNumber  *string `json:"smsPhoneNumber" binding:"omitempty,max=20"`
		IsCorporate     bool    `json:"isCorporate"`
		RoleName        string  `json:"roleName" binding:"required,roleoneof"`
		Status          string  `json:"status" binding:"statusoneof"`
		// User Details
		ClassId                    json.Number `json:"classId" binding:""`
		CountryOfResidenceIsoTwo   string      `json:"countryOfResidenceIsoTwo" binding:"max=2"`
		CountryOfCitizenshipIsoTwo string      `json:"countryOfCitizenshipIsoTwo" binding:"max=2"`
		DateOfBirth                *string     `json:"dateOfBirth"`
		DocumentType               string      `json:"documentType" binding:"documenttypeoneof"`
		DocumentPersonalId         string      `json:"documentPersonalId" binding:"max=255"`
		Fax                        string      `json:"fax" binding:"max=45"`
		HomePhoneNumber            string      `json:"homePhoneNumber" binding:"max=100"`
		InternalNotes              string      `json:"internalNotes"`
		OfficePhoneNumber          string      `json:"officePhoneNumber" binding:"max=100"`
		Position                   string      `json:"position" binding:"max=255"`
		UserGroupId                *uint64     `json:"userGroupId,omitempty"`
		CompanyDetails
		// Addresses
		PhysicalAdressValidator
		MailingAddressValidator
		BenificialOwnerValidator
	} `json:"data"`
	UserModel models.User `json:"-"`
}

type UserDetailsValidator struct {
	ClassId                    json.Number `json:"classId" binding:""`
	CountryOfResidenceIsoTwo   string      `json:"countryOfResidenceIsoTwo" binding:"max=2"`
	CountryOfCitizenshipIsoTwo string      `json:"countryOfCitizenshipIsoTwo" binding:"max=2"`
	DateOfBirth                *string     `json:"dateOfBirth"`
	DocumentType               *string     `json:"documentType" binding:"documenttypeoneof"`
	DocumentPersonalId         string      `json:"documentPersonalId" binding:"max=255"`
	Fax                        string      `json:"fax" binding:"max=45"`
	HomePhoneNumber            string      `json:"homePhoneNumber" binding:"max=100"`
	InternalNotes              string      `json:"internalNotes"`
	OfficePhoneNumber          string      `json:"officePhoneNumber" binding:"max=100"`
	Position                   string      `json:"position" binding:"max=255"`
	UserGroupId                *uint64     `json:"userGroupId,omitempty"`
}

// MailingAddressValidator is mailing adress validation model
type MailingAddressValidator struct {
	MaZipPostalCode   string `json:"maZipPostalCode" binding:"max=45"`
	MaStateProvRegion string `json:"maStateProvRegion" binding:"max=255"`
	MaPhoneNumber     string `json:"maPhoneNumber" binding:"max=255"`
	MaName            string `json:"maName" binding:"max=255"`
	MaCountryIso2     string `json:"maCountryISO2" binding:"max=2"`
	MaCity            string `json:"maCity" binding:"max=45"`
	MaAddress         string `json:"maAddress" binding:"max=255"`
	MaAddress2ndLine  string `json:"maAddress2ndLine" binding:"max=255"`
	MaAsPhysical      *bool  `json:"maAsPhysical"`
}

// BenificialOwnerValidator is beneficial owner validation model
type BenificialOwnerValidator struct {
	BoFullName           string `json:"boFullName" binding:"max=255"`
	BoPhoneNumber        string `json:"boPhoneNumber" binding:"max=255"`
	BoDateOfBirthYear    uint64 `json:"boDateOfBirthYear"`
	BoDateOfBirthMonth   uint64 `json:"boDateOfBirthMonth"`
	BoDateOfBirthDay     uint64 `json:"boDateOfBirthDay"`
	BoDocumentPersonalId string `json:"boDocumentPersonalId" binding:"max=255"`
	BoDocumentType       string `json:"boDocumentType" binding:"documenttypeoneof"`
	BoAddress            string `json:"boAddress" binding:"max=255"`
	BoRealationship      string `json:"boRelationship" binding:"max=255"`
}

// PhysicalAdressValidator is physical adress validation model
type PhysicalAdressValidator struct {
	PaZipPostalCode   string `json:"paZipPostalCode" binding:"max=45"`
	PaAddress         string `json:"paAddress" binding:"max=255"`
	PaAddress2ndLine  string `json:"paAddress2ndLine" binding:"max=255"`
	PaCity            string `json:"paCity" binding:"max=45"`
	PaCountryIso2     string `json:"paCountryISO2" binding:"max=2"`
	PaStateProvRegion string `json:"paStateProvRegion" binding:"max=255"`
}

type CompanyDetails struct {
	CompanyName       string `json:"companyName" binding:"max=255"`
	CompanyType       string `json:"companyType" binding:"max=255"`
	CompanyRole       string `json:"companyRole" binding:"max=255"`
	DirectorFirstName string `json:"directorFirstName" binding:"max=255"`
	DirectorLastName  string `json:"directorLastName" binding:"max=255"`
}

// BindJSON binding from JSON
func (s *CreateUserValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.UserModel.Username = s.Data.Username
	s.UserModel.Email = s.Data.Email
	s.UserModel.FirstName = s.Data.FirstName
	s.UserModel.LastName = s.Data.LastName
	s.UserModel.MiddleName = s.Data.MiddleName
	s.UserModel.Nickname = s.Data.Nickname

	s.UserModel.Password = s.Data.Password
	// s.UserModel.GeneratePassword()

	s.UserModel.PhoneNumber = s.Data.PhoneNumber
	s.UserModel.SmsPhoneNumber = s.Data.SmsPhoneNumber
	s.UserModel.IsCorporate = &s.Data.IsCorporate
	s.UserModel.RoleName = s.Data.RoleName
	s.UserModel.UserGroupId = s.Data.UserGroupId

	s.UserModel.Status = s.Data.Status
	s.UserModel.Position = s.Data.Position
	s.UserModel.InternalNotes = s.Data.InternalNotes
	s.UserModel.ClassId = s.Data.ClassId
	s.UserModel.LastActedAt = time.Now()

	//s.UserModel.DateOfBirth = s.User.DateOfBirth // TODO: add validation
	if len(s.Data.DocumentType) > 0 {
		s.UserModel.DocumentType = &s.Data.DocumentType
	} else {
		s.UserModel.DocumentType = nil
	}

	s.UserModel.DocumentPersonalId = s.Data.DocumentPersonalId
	s.UserModel.CountryOfResidenceIsoTwo = s.Data.CountryOfResidenceIsoTwo
	s.UserModel.CountryOfCitizenshipIsoTwo = s.Data.CountryOfCitizenshipIsoTwo

	s.UserModel.HomePhoneNumber = s.Data.HomePhoneNumber
	s.UserModel.OfficePhoneNumber = s.Data.OfficePhoneNumber
	s.UserModel.Fax = s.Data.Fax

	// Company details
	s.UserModel.CompanyDetails.CompanyName = s.Data.CompanyName
	s.UserModel.CompanyDetails.CompanyType = s.Data.CompanyType
	s.UserModel.CompanyDetails.CompanyRole = s.Data.CompanyRole
	s.UserModel.CompanyDetails.DirectorFirstName = s.Data.DirectorFirstName
	s.UserModel.CompanyDetails.DirectorLastName = s.Data.DirectorLastName

	return nil
}
