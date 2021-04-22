package states

import (
	"errors"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type ProgressState struct {
	context *models.Verification
}

func NewProgressState(context *models.Verification) *ProgressState {
	return &ProgressState{context}
}

func (s ProgressState) HandleVerificationRequest() error {
	return errors.New("could not create verification request with status 'progress'")
}

func (s ProgressState) HandleAdminApprove() error {
	s.context.Status = models.VerificationStatusApproved
	return nil
}

func (s ProgressState) HandleAdminCancellation() error {
	s.context.Status = models.VerificationStatusCancelled
	return nil
}
