package adminprofilesrow

import (
	"github.com/Confialink/wallet-pkg-utils/timefmt"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

// AdminProfileRowBuilder builder for admin profiles
type AdminProfileRowBuilder struct {
	user         *models.User
	timeSettings *syssettings.TimeSettings
}

// NewRowBuilder returns admAdminRowBuilderinRowBuilder
func NewRowBuilder(user *models.User, timeSettings *syssettings.TimeSettings) *AdminProfileRowBuilder {
	return &AdminProfileRowBuilder{user, timeSettings}
}

// Call returns array of fields for admin row in cvs file
func (b *AdminProfileRowBuilder) Call() []string {
	return []string{
		b.profileType(),
		b.firstName(),
		b.lastName(),
		b.username(),
		b.email(),
		b.created(),
		b.position(),
		b.status(),
		b.phoneNumber(),
		b.class(),
	}
}

func (b *AdminProfileRowBuilder) profileType() string {
	return b.user.GetProfileType()
}

func (b *AdminProfileRowBuilder) firstName() string {
	return b.user.FirstName
}

func (b *AdminProfileRowBuilder) lastName() string {
	return b.user.LastName
}

func (b *AdminProfileRowBuilder) username() string {
	return b.user.Username
}

func (b *AdminProfileRowBuilder) email() string {
	return b.user.Email
}

func (b *AdminProfileRowBuilder) created() string {
	return timefmt.Format(b.user.CreatedAt, b.timeSettings.DateTimeFormat, b.timeSettings.Timezone)
}

func (b *AdminProfileRowBuilder) position() string {
	return b.user.Position
}

func (b *AdminProfileRowBuilder) status() string {
	return b.user.Status
}

func (b *AdminProfileRowBuilder) phoneNumber() string {
	return b.user.PhoneNumber
}

func (b *AdminProfileRowBuilder) class() string {
	if b.user.PermissionGroup == nil {
		return ""
	}
	return b.user.PermissionGroup.Name
}

func (b *AdminProfileRowBuilder) internalNotes() string {
	return b.user.InternalNotes
}
