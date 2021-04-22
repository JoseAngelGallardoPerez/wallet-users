package users

import (
	"math/rand"
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/jinzhu/gorm"
)

type ConfirmationCode struct {
	confirmationCodeRepository *repositories.ConfirmationCodeRepository
}

const (
	defaultCodeLength = 42
	codeSymbolStart   = 33
	codeSymbolEnd     = 126

	codeLength  = 6
	expiryTime  = 60
	codeCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func NewConfirmationCode(confirmationCodeRepository *repositories.ConfirmationCodeRepository) *ConfirmationCode {
	return &ConfirmationCode{
		confirmationCodeRepository: confirmationCodeRepository,
	}
}

func (c *ConfirmationCode) GeneratePasswordRecoveryCode(user *models.User) (*models.ConfirmationCode, error) {
	model := c.createVerificationCode(user, models.ConfirmationCodeSubjectPasswordRecovery, time.Duration(48)*time.Hour)
	return c.confirmationCodeRepository.CreateNewVerificationCode(model)
}

func (c *ConfirmationCode) GenerateSetPasswordCode(user *models.User) (*models.ConfirmationCode, error) {
	return c.GenerateCode(user, models.ConfirmationCodeSubjectSetPassword, 64, time.Duration(24)*time.Hour)
}

func (c *ConfirmationCode) GenerateCode(user *models.User, subject string, codeLength uint8, expires time.Duration) (*models.ConfirmationCode, error) {
	model := &models.ConfirmationCode{
		Code:      c.generateUniqCode(codeLength),
		Subject:   subject,
		UserUID:   user.UID,
		ExpiresAt: time.Now().Add(expires),
	}

	return c.confirmationCodeRepository.Create(model)
}

func (c *ConfirmationCode) FindByCodeAndSubject(code string, subject string) (*models.ConfirmationCode, error) {
	return c.confirmationCodeRepository.FindByCodeAndSubject(code, subject)
}

func (c *ConfirmationCode) FindByCodeAndUsername(code string, subject string, username string) (*models.ConfirmationCode, error) {
	return c.confirmationCodeRepository.FindByCodeSubjectAndUsername(code, subject, username)
}

func (c *ConfirmationCode) FindByCodeSubjectAndEmailOrUsername(code string, subject string, emailOrUsername string) (*models.ConfirmationCode, error) {
	return c.confirmationCodeRepository.FindByCodeSubjectAndEmailOrUsername(code, subject, emailOrUsername)
}

func (c *ConfirmationCode) DeleteConfirmationCode(code *models.ConfirmationCode) error {
	return c.confirmationCodeRepository.Delete(code)
}

func (c ConfirmationCode) WrapContext(tx *gorm.DB) *ConfirmationCode {
	c.confirmationCodeRepository = c.confirmationCodeRepository.WrapContext(tx)
	return &c
}

func (c *ConfirmationCode) generateUniqCode(length uint8) string {
	for {
		code := c.generateCode(length)
		_, err := c.confirmationCodeRepository.FindByCode(code)
		if err == gorm.ErrRecordNotFound {
			return code
		}
	}
}

func (c *ConfirmationCode) generateCode(length uint8) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = rune(codeSymbolStart + rand.Intn(codeSymbolEnd-codeSymbolStart+1))
	}

	return string(b)
}

func (c *ConfirmationCode) generateNewCode(length uint8) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = codeCharset[seededRand.Intn(len(codeCharset))]
	}
	return string(b)
}

func (c *ConfirmationCode) createVerificationCode(user *models.User, subject string, expires time.Duration) *models.ConfirmationCode {
	model := &models.ConfirmationCode{
		Code:      c.generateNewCode(codeLength),
		Subject:   subject,
		UserUID:   user.UID,
		ExpiresAt: time.Now().Add(expires),
	}
	return model
}

func (c *ConfirmationCode) CreateNewVerificationCode(user *models.User, subject string) (*models.ConfirmationCode, error) {
	model := c.createVerificationCode(user, subject, expiryTime*time.Minute)
	return c.confirmationCodeRepository.CreateNewVerificationCode(model)
}

func (c *ConfirmationCode) CheckPhoneCode(phoneCode string, user *models.User) error {
	return c.confirmationCodeRepository.CheckPhoneCode(phoneCode, user)
}

func (c *ConfirmationCode) CheckEmailCode(emailCode string, user *models.User) error {
	return c.confirmationCodeRepository.CheckEmailCode(emailCode, user)
}
