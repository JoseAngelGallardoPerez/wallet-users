package userprofilesrow

import (
	"github.com/Confialink/wallet-pkg-utils/timefmt"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

type userFields struct {
	user         *models.User
	timeSettings *syssettings.TimeSettings
}

func newUserFields(user *models.User, timeSettings *syssettings.TimeSettings) *userFields {
	return &userFields{user, timeSettings}
}

func (f *userFields) call() []string {
	return []string{
		f.profileType(),
		f.firstName(),
		f.lastName(),
		f.username(),
		f.email(),
		f.created(),
		f.comanyName(),
		f.status(),
		f.dateOfBirth(),
		f.documnetType(),
		f.documnetNumber(),
		f.countryOfResidence(),
		f.countryOfCitizenship(),
		f.smsPhoneNumber(),
		f.phoneNumber(),
		f.homePhone(),
		f.officePhone(),
		f.fax(),
		f.group(),
	}
}

func (f *userFields) profileType() string {
	return f.user.GetProfileType()
}

func (f *userFields) firstName() string {
	return f.user.FirstName
}

func (f *userFields) lastName() string {
	return f.user.LastName
}

func (f *userFields) username() string {
	return f.user.Username
}

func (f *userFields) email() string {
	return f.user.Email
}

func (f *userFields) created() string {
	return timefmt.Format(f.user.CreatedAt, f.timeSettings.DateTimeFormat, f.timeSettings.Timezone)
}

func (f *userFields) comanyName() string {
	return f.user.CompanyDetails.CompanyName
}

func (f *userFields) status() string {
	return f.user.Status
}

func (f *userFields) dateOfBirth() string {
	if f.user.DateOfBirth == nil {
		return ""
	}
	return f.user.DateOfBirth.String() // TODO: check format
}

func (f *userFields) documnetType() string {
	if f.user.UserDetails.DocumentType != nil {
		return *f.user.UserDetails.DocumentType
	}
	return ""
}

func (f *userFields) documnetNumber() string {
	return f.user.UserDetails.DocumentPersonalId
}

func (f *userFields) countryOfResidence() string {
	return f.user.UserDetails.CountryOfResidenceIsoTwo
}

func (f *userFields) countryOfCitizenship() string {
	return f.user.UserDetails.CountryOfCitizenshipIsoTwo
}

// TODO: add sms phone number
func (f *userFields) smsPhoneNumber() string {
	return ""
}

func (f *userFields) phoneNumber() string {
	return f.user.PhoneNumber
}

func (f *userFields) homePhone() string {
	return f.user.UserDetails.HomePhoneNumber
}

func (f *userFields) officePhone() string {
	return f.user.UserDetails.OfficePhoneNumber
}

func (f *userFields) fax() string {
	return f.user.UserDetails.Fax
}

func (f *userFields) group() string {
	if f.user.UserGroup != nil {
		return f.user.UserGroup.Name
	}
	return ""
}
