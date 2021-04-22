package verification

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
)

const maxTypeFiles = 6

type Validator struct {
	verificationsRepository     *repositories.VerificationRepository
	verificationFilesRepository *repositories.VerificationFilesRepository
}

func NewValidator(verificationsRepository *repositories.VerificationRepository,
	verificationFilesRepository *repositories.VerificationFilesRepository) *Validator {
	return &Validator{verificationsRepository, verificationFilesRepository}
}

func (v *Validator) CanCreate(verificationType, userId string) errors.TypedError {
	verification, err := v.verificationsRepository.FindByUIDAndType(userId, verificationType)
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != nil {
		return &errors.PrivateError{OriginalError: err}
	}

	return v.validateCount(verification.ID)
}

func (v *Validator) validateCount(verificationID uint32) errors.TypedError {
	count, err := v.verificationFilesRepository.CountByVerificationID(verificationID)
	if err != nil {
		return &errors.PrivateError{OriginalError: err}
	}

	if count >= maxTypeFiles {
		return responses.GetTypedError(responses.MaxVerificationFiles)
	}
	return nil
}
