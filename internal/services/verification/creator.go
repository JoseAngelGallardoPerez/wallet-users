package verification

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
)

type Creator struct {
	verificationsRepository     *repositories.VerificationRepository
	verificationFilesRepository *repositories.VerificationFilesRepository
	db                          *gorm.DB
}

func NewCreator(verificationsRepository *repositories.VerificationRepository,
	verificationFilesRepository *repositories.VerificationFilesRepository,
	db *gorm.DB) *Creator {
	return &Creator{
		verificationsRepository,
		verificationFilesRepository,
		db,
	}
}

func (c *Creator) Call(verificationType, userId string, fileId uint64) (*models.Verification, errors.TypedError) {
	tx := c.beginTransaction()
	model, err := c.call(verificationType, userId, fileId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return model, nil
}

func (c *Creator) call(verificationType, userId string, fileId uint64) (*models.Verification, errors.TypedError) {
	model, err := c.verificationsRepository.FindByUIDAndType(userId, verificationType)
	if err == gorm.ErrRecordNotFound {
		return c.createWithFile(verificationType, userId, fileId)
	}
	if err != nil {
		return nil, &errors.PrivateError{OriginalError: err}
	}

	return c.addFile(model, fileId)
}

func (c *Creator) addFile(model *models.Verification, fileId uint64) (*models.Verification, errors.TypedError) {
	if typedErr := c.createFile(model.ID, fileId); typedErr != nil {
		return nil, typedErr
	}

	files, err := c.verificationFilesRepository.FindByVerificationID(model.ID)
	if err != nil {
		return nil, &errors.PrivateError{OriginalError: err}
	}
	model.Files = files
	return model, nil
}

func (c *Creator) createFile(verificationId uint32, fileId uint64) errors.TypedError {
	file := models.VerificationFile{
		VerificationID: verificationId,
		FileID:         fileId,
	}
	if _, err := c.verificationFilesRepository.Create(&file); err != nil {
		return &errors.PrivateError{OriginalError: err}
	}
	return nil
}

func (c *Creator) createWithFile(verificationType, userId string, fileId uint64) (*models.Verification, errors.TypedError) {
	model := models.Verification{
		Status:  models.VerificationStatusPending,
		Type:    verificationType,
		UserUID: userId,
	}

	created, err := c.verificationsRepository.Create(&model)
	if err != nil {
		return nil, &errors.PrivateError{OriginalError: err}
	}

	return c.addFile(created, fileId)
}

func (c *Creator) beginTransaction() *gorm.DB {
	tx := c.db.Begin()
	c.verificationsRepository = c.verificationsRepository.WrapContext(tx)
	c.verificationFilesRepository = c.verificationFilesRepository.WrapContext(tx)
	return tx
}
