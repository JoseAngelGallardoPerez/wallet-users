package states

import (
	"errors"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type ApprovedState struct {
	context *models.Verification
}

func NewApprovedState(context *models.Verification) *ApprovedState {
	return &ApprovedState{context}
}

func (s ApprovedState) HandleVerificationRequest() error {
	return errors.New("could not create verification request with status 'approved'")
}

func (s ApprovedState) HandleAdminApprove() error {
	return errors.New("could not approve verification with status 'approved'")
}

func (s ApprovedState) HandleAdminCancellation() error {
	return errors.New("could not cancel verification with status 'approved'")
}
