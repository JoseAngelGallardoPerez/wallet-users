package verificationstates

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/verification_states/states"
)

type VerificationStater interface {
	HandleVerificationRequest() error
	HandleAdminApprove() error
	HandleAdminCancellation() error
}

func NewVerificationState(v *models.Verification) VerificationStater {
	var state VerificationStater
	switch v.Status {
	case models.VerificationStatusPending:
		state = states.NewPendingState(v)
	case models.VerificationStatusProgress:
		state = states.NewProgressState(v)
	case models.VerificationStatusApproved:
		state = states.NewApprovedState(v)
	case models.VerificationStatusCancelled:
		state = states.NewCancelledState(v)
	}
	return state
}
