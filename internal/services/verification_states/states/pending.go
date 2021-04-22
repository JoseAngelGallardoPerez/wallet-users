package states

import (
	"github.com/Confialink/wallet-users/internal/db/models"
)

type PendingState struct {
	context *models.Verification
}

func NewPendingState(context *models.Verification) *PendingState {
	return &PendingState{context}
}

func (s PendingState) HandleVerificationRequest() error {
	s.context.Status = models.VerificationStatusProgress
	return nil
}

func (s PendingState) HandleAdminApprove() error {
	s.context.Status = models.VerificationStatusApproved
	return nil
}

func (s PendingState) HandleAdminCancellation() error {
	s.context.Status = models.VerificationStatusCancelled
	return nil
}
