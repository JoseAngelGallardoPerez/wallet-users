package users

import (
	"strings"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/Confialink/wallet-pkg-utils/pointer"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/files"
	"github.com/Confialink/wallet-users/internal/services/gdpr"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

type UserService struct {
	db                      *gorm.DB
	userRepository          *repositories.UsersRepository
	companyService          *CompanyService
	addressRepo             *repositories.AddressRepository
	attributeService        *AttributeService
	logger                  log15.Logger
	passwordService         *services.Password
	notificationsService    *notifications.Notifications
	confirmationCodeService *ConfirmationCode
	settings                *syssettings.SysSettings
	files                   *files.FilesService
	gdprService             *gdpr.Service
}

func NewUserService(
	db *gorm.DB,
	userRepository *repositories.UsersRepository,
	companyService *CompanyService,
	addressRepo *repositories.AddressRepository,
	attributeService *AttributeService,
	logger log15.Logger,
	pwdService *services.Password,
	notificationsService *notifications.Notifications,
	confirmationCodeService *ConfirmationCode,
	settings *syssettings.SysSettings,
	files *files.FilesService,
	gdprService *gdpr.Service,
) *UserService {
	return &UserService{
		db,
		userRepository,
		companyService,
		addressRepo,
		attributeService,
		logger,
		pwdService,
		notificationsService,
		confirmationCodeService,
		settings,
		files,
		gdprService,
	}
}

// Updating in DB user, Mailing Addresses, Physical Addresses and Attributes
func (this *UserService) Update(user *models.User, tx *gorm.DB) error {
	var userRepo *repositories.UsersRepository

	var localTransaction bool
	if tx == nil {
		localTransaction = true
		tx = this.db.Begin()
	}
	userRepo = this.userRepository.WrapContext(tx)

	if user.UserGroupId == nil {
		user.UserGroup = nil
	}

	// Update Addresses
	addressRepo := this.addressRepo.WrapContext(tx)
	if err := this.attachAddresses(user.MailingAddresses, models.AddressTypeMailing, user.UID, addressRepo); err != nil {
		if localTransaction {
			tx.Rollback()
		}
		return err
	}
	if err := this.attachAddresses(user.PhysicalAddresses, models.AddressTypePhysical, user.UID, addressRepo); err != nil {
		if localTransaction {
			tx.Rollback()
		}
		return err
	}

	// Update Attributes
	// TODO: check!!! Perhaps it duplicates attributes
	if err := this.attributeService.AttachAttributes(user.Attributes, user.UID, tx); err != nil {
		if localTransaction {
			tx.Rollback()
		}
		return err
	}

	// Update Company Details
	if err := this.companyService.UpdateCompanyDetails(user, tx); err != nil {
		if localTransaction {
			tx.Rollback()
		}
		return err
	}

	// Use Save() to save empty fields
	_, err := userRepo.Save(user)
	if err != nil {
		this.logger.Error("failed update user", "error", err)
		if localTransaction {
			tx.Rollback()
		}
		return err
	}

	if localTransaction {
		tx.Commit()
	}
	return nil
}

func (this *UserService) Create(initUser *models.User, confirmed bool, newPasswordRequired bool, tx *gorm.DB) (*models.User, error) {
	var userRepo *repositories.UsersRepository
	var localTransaction bool
	if tx == nil {
		localTransaction = true
		tx = this.db.Begin()
	}
	userRepo = this.userRepository.WrapContext(tx)

	if !initUser.IsPasswordEncrypted() {
		// Generates a hashed version of our password
		hash, err := this.passwordService.UserHashPassword(initUser.Password)
		if err != nil {
			this.logger.Error("failed create hash password", "error", err)
			if localTransaction {
				tx.Rollback()
			}
			return nil, err
		}

		initUser.Password = hash
	}

	if newPasswordRequired {
		initUser.ChallengeName = pointer.ToString(models.ChallengeNameNewPasswordRequired)
	}

	user, err := userRepo.Create(initUser, confirmed)
	if err != nil {
		this.logger.Error("failed create user", "error", err)
		if localTransaction {
			tx.Rollback()
		}
		return nil, err
	}

	// Addresses
	addressRepo := this.addressRepo.WrapContext(tx)
	if err := this.attachAddresses(initUser.MailingAddresses, models.AddressTypeMailing, user.UID, addressRepo); err != nil {
		if localTransaction {
			tx.Rollback()
		}
		return nil, err
	}
	if err := this.attachAddresses(initUser.PhysicalAddresses, models.AddressTypePhysical, user.UID, addressRepo); err != nil {
		if localTransaction {
			tx.Rollback()
		}
		return nil, err
	}

	if err := this.attributeService.AttachAttributes(initUser.Attributes, user.UID, tx); err != nil {
		if localTransaction {
			tx.Rollback()
		}
		return nil, err
	}

	if localTransaction {
		tx.Commit()
	}

	return user, nil
}

func (this *UserService) CreateNew(user *models.User) (*models.User, error) {
	tx := this.db.Begin()
	user, err := this.Create(user, true, false, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// TODO: Move to the Message broker
	gdprSettings, err := this.settings.GetGDPRSettings()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if gdprSettings.Enabled {
		gdprBytes, err := this.gdprService.GdprHtmlBytes()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = this.files.Upload(gdprBytes, "gdpr-accepted.pdf", user.UID, true, true, "gdpr")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return user, nil
}

func (this *UserService) ForgotPassword(emailOrPhoneNumber string) error {
	user, err := this.userRepository.FindByEmailOrPhoneNumber(emailOrPhoneNumber)
	if err != nil {
		var vErrs []errors.ValidationError
		vErrs = append(vErrs, errors.ValidationError{
			Title:  "Sorry, unrecognized e-mail address or phone number.",
			Source: "email",
			Code:   responses.UnknownEmailOrPhoneNumber,
		})
		this.logger.Error("can not find user with email or phone number", "error", err)
		return &errors.ValidationErrors{Errors: vErrs}
	}

	if user.IsBlocked() {
		return &errors.PublicError{
			Code: responses.ActionCannotBePerformed,
		}
	}

	if !user.IsActive() {
		return &errors.PublicError{
			Title: "User is not active",
			Code:  responses.CodeResetPasswordIsNotAllowed,
		}
	}

	codeModel, err := this.confirmationCodeService.GeneratePasswordRecoveryCode(user)
	if err != nil {
		return err
	}

	// If a user enters a phone number - we send sms. Otherwise we send an email.
	checkEmail := strings.Index(emailOrPhoneNumber, "@")
	methods := []string{"sms"}
	if checkEmail > -1 {
		methods = []string{"email"}
	}

	if _, err = this.notificationsService.PasswordRecovery(user.UID, codeModel.Code, methods); err != nil {
		this.logger.Error("сan not send notification", "error", err)
		return err
	}

	return nil
}

func (this *UserService) ConfirmForgotPassword(password string, code string, emailOrUsername string) error {
	codeModel, err := this.confirmationCodeService.FindByCodeSubjectAndEmailOrUsername(code, models.ConfirmationCodeSubjectPasswordRecovery, emailOrUsername)
	if err != nil {
		return err
	}

	user := codeModel.User
	hash, err := this.passwordService.UserHashPassword(password)
	user.Password = hash

	err = this.confirmationCodeService.DeleteConfirmationCode(codeModel)
	if err != nil {
		return err
	}

	_, err = this.userRepository.Update(user)
	return err
}

func (this *UserService) CheckPhoneCode(phoneCode string, user *models.User) (*models.User, error) {
	err := this.confirmationCodeService.CheckPhoneCode(phoneCode, user)
	if err != nil {
		return nil, err
	}

	user.IsPhoneConfirmed = true

	return this.userRepository.Update(user)
}

func (this *UserService) CheckEmailCode(emailCode string, user *models.User) (*models.User, error) {
	err := this.confirmationCodeService.CheckEmailCode(emailCode, user)
	if err != nil {
		return nil, err
	}

	user.IsEmailConfirmed = true

	return this.userRepository.Update(user)
}

func (s *UserService) ResetPassword(newPassword, code string) error {
	codeModel, err := s.confirmationCodeService.FindByCodeAndSubject(code, models.ConfirmationCodeSubjectPasswordRecovery)
	if err != nil {
		return err
	}

	err = s.SetPassword(codeModel.User, newPassword)
	if err != nil {
		return err
	}

	err = s.confirmationCodeService.DeleteConfirmationCode(codeModel)
	if err != nil {
		return err
	}
	return nil
}

// Set new password for user
func (s *UserService) SetPassword(user *models.User, password string) error {
	hash, err := s.passwordService.UserHashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hash
	_, err = s.userRepository.Update(user)
	return err
}

// Attaches addresses to the user
func (s *UserService) attachAddresses(addresses []*models.Address, addressType, userId string, repo *repositories.AddressRepository) error {
	if len(addresses) > 0 {
		for _, address := range addresses {
			address.UserID = userId
			address.Type = addressType
			if err := repo.Save(address); err != nil {
				return err
			}
		}
	} else {
		if err := repo.DeleteByTypeForUser(userId, addressType); err != nil {
			return err
		}
	}

	return nil
}
