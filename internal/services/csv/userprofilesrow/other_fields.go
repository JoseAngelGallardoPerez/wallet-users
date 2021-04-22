package userprofilesrow

import (
	"github.com/Confialink/wallet-users/internal/db/models"
)

type otherFields struct {
	userDetails *models.UserDetails
}

func newOtherFields(userDetails *models.UserDetails) *otherFields {
	return &otherFields{userDetails}
}

func (f *otherFields) call() []string {
	return []string{
		f.userDetails.InternalNotes,
	}
}
