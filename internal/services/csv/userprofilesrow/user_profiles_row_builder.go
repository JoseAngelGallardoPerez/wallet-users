package userprofilesrow

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

// UserProfileRowBuilder service to build one row for user profile csv
type UserProfileRowBuilder struct {
	user         *models.User
	timeSettings *syssettings.TimeSettings
}

// NewUserProfileRowBuilder returns new UserProfileRowBuilder
func NewUserProfileRowBuilder(user *models.User, timeSettings *syssettings.TimeSettings) *UserProfileRowBuilder {
	return &UserProfileRowBuilder{user, timeSettings}
}

// Call returns one row for user profile
func (b *UserProfileRowBuilder) Call() []string {
	userFields := newUserFields(b.user, b.timeSettings).call()

	//physicalAddressFields := newPhysicalAddress(&b.user.PhysicalAddress).call()
	//mailingAddressFields := newMailingAddress(&b.user.MailingAddress).call()

	otherFields := newOtherFields(&b.user.UserDetails).call()

	result := make([]string, 0)
	result = append(result, userFields...)
	//result = append(result, physicalAddressFields...)
	//result = append(result, mailingAddressFields...)
	result = append(result, otherFields...)
	return result
}
