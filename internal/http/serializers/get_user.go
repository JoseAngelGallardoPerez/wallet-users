package serializers

import (
	model_serializer "github.com/Confialink/wallet-pkg-model_serializer"

	"github.com/Confialink/wallet-users/internal/db/models"
)

var getUserFields = []interface{}{
	"UID", "Email", "Username", "FirstName", "LastName", "MiddleName", "Nickname", "PhoneNumber", "SmsPhoneNumber",
	"CompanyID", "ProfileImageID",
	map[string][]interface{}{"CompanyDetails": {"ID", "CompanyName", "CompanyType", "CompanyRole", "DirectorFirstName", "DirectorLastName"}},
	"IsCorporate", "RoleName", "ParentId", "Status",
	map[string][]interface{}{"PhysicalAddresses": {"ID", "CountryIsoTwo", "Region", "City", "ZipCode", "Address", "AddressSecondLine", "Name", "PhoneNumber", "Description", "Latitude", "Longitude"}},
	map[string][]interface{}{"MailingAddresses": {"ID", "CountryIsoTwo", "Region", "City", "ZipCode", "Address", "AddressSecondLine", "Name", "PhoneNumber", "Description", "Latitude", "Longitude"}},
	map[string][]interface{}{"UserGroup": {"ID", "Name", "Description"}},
	"UserGroupId", "CreatedAt", "UpdatedAt", "LastLoginAt", "LastLoginIp", "ChallengeName", "ClassId",
	"CountryOfResidenceIsoTwo", "CountryOfCitizenshipIsoTwo", "DateOfBirth",
	"DocumentType", "DocumentPersonalId", "Fax", "HomePhoneNumber", "InternalNotes", "OfficePhoneNumber", "Position",
	"BlockedUntil", "LastActedAt", "Attributes",
}

type getUser struct {
	user *models.User
}

func NewGetUser(user *models.User) *getUser {
	return &getUser{user}
}

func (s *getUser) Serialize() map[string]interface{} {
	return model_serializer.Serialize(s.user, getUserFields)
}
