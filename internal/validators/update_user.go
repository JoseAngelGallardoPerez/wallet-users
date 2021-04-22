package validators

import (
	"encoding/json"

	"github.com/Confialink/wallet-pkg-utils/pointer"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/Confialink/wallet-users/internal/db/models"
)

// UpdateUserValidator containing two parts:
// - Validator: write the form/json checking rule according to the doc https://github.com/go-playground/validator
// - DataModel: fill with data from Validator after invoking ctx.ShouldBindWith(self)
// Then, you can just call CreateUserValidator.userModel after the data is ready in DataModel.
type UpdateUserValidator struct {
	Data struct {
		UID            string  `json:"-"`
		Username       string  `json:"username" binding:"omitempty,min=4,max=50,usernameChars"`
		Email          string  `json:"email" binding:"email"`
		FirstName      string  `json:"firstName" binding:"max=255"`
		LastName       string  `json:"lastName" binding:"max=255"`
		MiddleName     string  `json:"middleName" binding:"max=255"`
		Nickname       string  `json:"nickname" binding:"max=100"`
		PhoneNumber    string  `json:"phoneNumber" binding:"max=100"`
		SmsPhoneNumber *string `json:"smsPhoneNumber" binding:"omitempty,max=20"`
		IsCorporate    bool    `json:"isCorporate"`
		RoleName       string  `json:"roleName" binding:"required,roleoneof"`
		Status         string  `json:"status" binding:"statusoneof"`

		CompanyID         *uint64 `json:"companyID"`
		CompanyName       string  `json:"companyName" binding:"max=255"`
		CompanyType       string  `json:"companyType" binding:"max=255"`
		CompanyRole       string  `json:"companyRole" binding:"max=255"`
		DirectorFirstName string  `json:"directorFirstName" binding:"max=255"`
		DirectorLastName  string  `json:"directorLastName" binding:"max=255"`

		// Password
		PreviousPassword *string `json:"previousPassword,omitempty"`
		Password         *string `json:"password,omitempty"`
		ConfirmPassword  *string `json:"confirmPassword,omitempty"`

		// User Details
		ClassId                    json.Number `json:"classId" binding:"required"`
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

		// Addresses
		PhysicalAdressValidator
		MailingAddressValidator
		BenificialOwnerValidator
	} `json:"data"`
	UserModel   models.User `json:"-"`
	CurrentUser *userpb.User
}

// GetUpdateUserValidator returns UpdateUserValidator
func GetUpdateUserValidator() UpdateUserValidator {
	validator := UpdateUserValidator{}
	return validator
}

// GetUpdateUserValidatorFillWith ...
func GetUpdateUserValidatorFillWith(user models.User) UpdateUserValidator {

	v := GetUpdateUserValidator()
	v.Data.UID = user.UID
	v.Data.Username = user.Username
	v.Data.Email = user.Email
	v.Data.FirstName = user.FirstName
	v.Data.LastName = user.LastName
	v.Data.MiddleName = user.MiddleName
	v.Data.Nickname = user.Nickname
	v.Data.PhoneNumber = user.PhoneNumber
	v.Data.SmsPhoneNumber = user.SmsPhoneNumber
	v.Data.IsCorporate = *user.IsCorporate
	v.Data.RoleName = user.RoleName
	v.Data.Status = user.Status
	v.Data.Position = user.Position
	v.Data.InternalNotes = user.InternalNotes
	v.Data.ClassId = user.ClassId

	v.Data.CompanyID = user.CompanyID
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

// BindJSON binding from JSON
func (s *UpdateUserValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.UserModel.Username = s.Data.Username

	s.UserModel.UID = s.Data.UID
	s.UserModel.Email = s.Data.Email
	s.UserModel.FirstName = s.Data.FirstName
	s.UserModel.LastName = s.Data.LastName
	s.UserModel.MiddleName = s.Data.MiddleName
	s.UserModel.Nickname = s.Data.Nickname
	s.UserModel.PhoneNumber = s.Data.PhoneNumber
	s.UserModel.SmsPhoneNumber = s.Data.SmsPhoneNumber
	s.UserModel.IsCorporate = &s.Data.IsCorporate
	s.UserModel.Position = s.Data.Position

	s.UserModel.CompanyID = s.Data.CompanyID
	s.UserModel.CompanyDetails.CompanyName = s.Data.CompanyName
	s.UserModel.CompanyDetails.CompanyType = s.Data.CompanyType
	s.UserModel.CompanyDetails.CompanyRole = s.Data.CompanyRole
	s.UserModel.CompanyDetails.DirectorFirstName = s.Data.DirectorFirstName
	s.UserModel.CompanyDetails.DirectorLastName = s.Data.DirectorLastName

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

	// only admin or root can change next fields
	currentUser := getCurrentUser(ctx)
	s.CurrentUser = currentUser
	if currentUser.RoleName == models.RoleRoot || currentUser.RoleName == models.RoleAdmin {
		s.UserModel.ClassId = s.Data.ClassId
		s.UserModel.Status = s.Data.Status
		s.UserModel.UserGroupId = s.Data.UserGroupId
		s.UserModel.UserGroup = nil // it is required to update UserGroupId
		s.UserModel.InternalNotes = s.Data.InternalNotes
	}

	if s.Data.Password != nil && len(*s.Data.Password) > 0 {
		s.UserModel.Password = *s.Data.Password
	}

	return nil
}

// getCurrentUser retrieve current user from gin context
func getCurrentUser(ctx *gin.Context) *userpb.User {
	user, ok := ctx.Get("_user")
	if !ok {
		return nil
	}
	return user.(*userpb.User)
}
