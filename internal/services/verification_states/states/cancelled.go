package states

import (
	"errors"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type CancelledState struct {
	context *models.Verification
}

func NewCancelledState(context *models.Verification) *CancelledState {
	return &CancelledState{context}
}

func (s CancelledState) HandleVerificationRequest() error {
	return errors.New("could not create verification request with status 'cancelled'")
}

func (s CancelledState) HandleAdminApprove() error {
	return errors.New("could not approve verification with status 'cancelled'")
}

func (s CancelledState) HandleAdminCancellation() error {
	return errors.New("could not cancel verification with status 'cancelled'")
}
