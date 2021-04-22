package repositories

import (
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"

	"github.com/jinzhu/gorm"
)

type ConfirmationCodeRepository struct {
	DB *gorm.DB
}

func NewConfirmationCodeRepository(DB *gorm.DB) *ConfirmationCodeRepository {
	return &ConfirmationCodeRepository{DB: DB}
}

func (repo *ConfirmationCodeRepository) FindByCodeAndSubject(code string, subject string) (*models.ConfirmationCode, error) {
	model := &models.ConfirmationCode{}
	if err := repo.DB.Where("confirmation_codes.code = ? AND confirmation_codes.subject = ?", code, subject).
		Preload("User").
		First(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (repo *ConfirmationCodeRepository) FindByCode(code string) (*models.ConfirmationCode, error) {
	model := &models.ConfirmationCode{}
	if err := repo.DB.Where("confirmation_codes.code = ?", code).
		First(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (repo *ConfirmationCodeRepository) FindByCodeSubjectAndUsername(code string, subject string, username string) (*models.ConfirmationCode, error) {
	model := &models.ConfirmationCode{}

	query := repo.DB.Joins("LEFT JOIN users ON users.uid = confirmation_codes.user_uid")

	if err := query.Where("confirmation_codes.code = ? AND confirmation_codes.subject = ? AND users.username = ?", code, subject, username).
		Preload("User").
		First(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (repo *ConfirmationCodeRepository) FindByCodeSubjectAndEmailOrUsername(code string, subject string, emailOrUsername string) (*models.ConfirmationCode, error) {
	model := &models.ConfirmationCode{}

	query := repo.DB.Joins("LEFT JOIN users ON users.uid = confirmation_codes.user_uid")

	if err := query.Where("confirmation_codes.code = ? AND confirmation_codes.subject = ? AND (users.email = ? OR users.username = ?)", code, subject, emailOrUsername, emailOrUsername).
		Preload("User").
		First(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (repo *ConfirmationCodeRepository) Create(code *models.ConfirmationCode) (*models.ConfirmationCode, error) {
	if err := repo.DB.Create(code).Error; err != nil {
		return nil, err
	}
	return code, nil
}

func (repo *ConfirmationCodeRepository) Delete(code *models.ConfirmationCode) error {
	if err := repo.DB.Delete(code).Error; err != nil {
		return err
	}
	return nil
}

func (copy ConfirmationCodeRepository) WrapContext(db *gorm.DB) *ConfirmationCodeRepository {
	copy.DB = db
	return &copy
}

func (repo *ConfirmationCodeRepository) CreateNewVerificationCode(code *models.ConfirmationCode) (*models.ConfirmationCode, error) {
	if err := repo.DB.Delete(models.ConfirmationCode{}, "user_uid = ? AND subject = ?", code.UserUID, code.Subject).Error; err != nil {
		return nil, err
	}

	return repo.Create(code)
}

func (repo *ConfirmationCodeRepository) CheckPhoneCode(phoneCode string, user *models.User) error {
	return repo.DB.
		Where("user_uid = ?", user.UID).
		Where("subject = ?", models.ConfirmationCodeSubjectPhoneVerificationCode).
		Where("code = ?", phoneCode).
		Where("expires_at >= ?", time.Now()).
		First(&models.ConfirmationCode{}).Error
}

func (repo *ConfirmationCodeRepository) CheckEmailCode(emailCode string, user *models.User) error {
	return repo.DB.
		Where("user_uid = ?", user.UID).
		Where("subject = ?", models.ConfirmationCodeSubjectEmailVerificationCode).
		Where("code = ?", emailCode).
		Where("expires_at >= ?", time.Now()).
		First(&models.ConfirmationCode{}).Error
}
