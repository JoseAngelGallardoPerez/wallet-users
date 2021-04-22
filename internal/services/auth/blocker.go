package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

type Blocker struct {
	blockedIpRepository       *repositories.BlockedIpsRepository
	failAuthAttemptRepository *repositories.FailAuthAttemptRepository
	userRepository            *repositories.UsersRepository
	logger                    log15.Logger
	loginSettings             *syssettings.LoginSecuritySettings
	notificationsService      *notifications.Notifications
	sysSettings               *syssettings.SysSettings
}

const (
	ipBlockingDurationMin   = 10
	userBlockingDurationMin = 30
)

func (b *Blocker) CheckIP(ip string) error {
	if len(ip) < 1 {
		return nil
	}
	blockedIP, err := b.blockedIpRepository.FindByIp(ip)
	if err != nil {
		return nil
	}

	if b.unblockIPIfNeed(blockedIP) {
		return nil
	}

	return errors.New("user IP is blocked")
}

func (b *Blocker) CheckUser(emailOrPhoneNumber string) *responses.Error {
	user, err := b.userRepository.FindByEmailOrPhoneNumber(emailOrPhoneNumber)
	if err != nil {
		return nil
	}

	if user.IsBlocked() {
		e := responses.NewCommonError().ApplyCode(responses.CodeUserIsBlocked)
		meta := struct {
			Username string `json:"username"`
		}{user.Username}
		e.Meta = meta
		e.Details = fmt.Sprintf("Sorry, user %s is blocked", user.Username)
		return e
	}

	return nil
}

func (b *Blocker) CheckPhoneConfirmed(emailOrPhoneNumber string) *responses.Error {
	user, err := b.userRepository.FindByEmailOrPhoneNumber(emailOrPhoneNumber)
	if err != nil {
		return nil
	}

	if !user.IsPhoneConfirmed {
		e := responses.NewCommonError().ApplyCode(responses.PhoneNumberIsNotConfirmed)
		meta := struct {
			Username    string `json:"username"`
			PhoneNumber string `json:"phoneNumber"`
		}{user.Username, user.PhoneNumber}
		e.Meta = meta
		e.Details = fmt.Sprintf("Sorry, phone number %s of user %s is not confirmed", user.PhoneNumber, user.Username)
		return e
	}

	return nil
}

func (b *Blocker) AddIPFailAttempt(ip string) bool {
	settings := b.loginSecuritySettings()
	if settings.FailedLoginUserUse == 0 {
		return false
	}

	if len(ip) > 1 && !b.blockIPIfNeed(ip) {
		b.createFailAttempt(ip, "")
		return true
	}

	return false
}

func (b *Blocker) AddUserFailAttempt(emailOrUsername, ip string) bool {
	user, err := b.userRepository.FindByEmailOrPhoneNumber(emailOrUsername)
	if err != nil {
		b.logger.Error("failed AddUserFailAttempt", "error", err)
		return false
	}

	settings := b.loginSecuritySettings()

	isNeedCreateUserFailAttempt := settings.FailedLoginUsernameUse && !b.blockUserIfNeed(user)
	isNeedCreateIPFailAttempt := len(ip) > 1 && !b.blockIPIfNeed(ip)
	if isNeedCreateUserFailAttempt || isNeedCreateIPFailAttempt {
		b.createFailAttempt(ip, user.UID)
		return true
	}
	return false
}

func (b *Blocker) ClearAllOldFailAttempts(user *models.User) {
	b.deleteOldAttempts(user.UID)
}

func (b *Blocker) LoadSettings() {
	var err error
	b.loginSettings, err = b.sysSettings.GetLoginSecuritySettings()
	if err != nil {
		b.logger.Error("failed fetch login security settings", "error", err)
	}
}

// return true if block IP
func (b *Blocker) blockIPIfNeed(ip string) bool {
	settings := b.loginSecuritySettings()
	if settings.FailedLoginUserUse == 0 {
		return false
	}

	timeFrom := b.createUntilTime(-ipBlockingDurationMin)
	countAttempts := b.failAuthAttemptRepository.GetCountByIp(ip, timeFrom)

	// we decrease 1 because we don't plus current attempt
	if countAttempts >= (uint32(settings.FailedLoginUserUse) - 1) {
		blockedUntil := b.createUntilTime(int64(settings.FailedLoginUserWindow))
		err := b.blockedIpRepository.Create(ip, blockedUntil)
		if err != nil {
			b.logger.Error("failed block ip", "error", err)
			return false
		}
		if err := b.failAuthAttemptRepository.DeleteAllByIp(ip); err != nil {
			b.logger.Error("failed delete failed auth attempts", "error", err)
		}

		return true
	}

	return false
}

// return true if unblock IP
func (b *Blocker) unblockIPIfNeed(blockedIP *models.BlockedIp) bool {
	if blockedIP.BlockedUntil.Before(time.Now()) {
		if err := b.blockedIpRepository.Delete(blockedIP); err != nil {
			b.logger.Error("failed delete blocked IP", "error", err)
			return false
		}
		return true
	}

	return false
}

// return true if block user
func (b *Blocker) blockUserIfNeed(user *models.User) bool {
	settings := b.loginSecuritySettings()
	if settings.FailedLoginUsernameLimit == 0 {
		return false
	}
	timeFrom := b.createUntilTime(-userBlockingDurationMin)
	countAttempts := b.failAuthAttemptRepository.GetCountByUID(user.UID, timeFrom)

	// we decrease 1 because we don't plus current attempt
	if countAttempts >= (uint32(settings.FailedLoginUsernameLimit) - 1) {
		blockedUntil := b.createUntilTime(int64(settings.FailedLoginUsernameCleanup))
		user.BlockedUntil = &blockedUntil
		err := b.userRepository.DB.Model(user).Updates(
			models.User{Status: user.Status,
				UserDetails: models.UserDetails{BlockedUntil: user.BlockedUntil}}).Error
		if err != nil {
			b.logger.Error("failed block user", "error", err)
		}
		//do not remove fail attempts, we use them for checking ip blocking
		if err := b.failAuthAttemptRepository.ImpersonateAllByUID(user.UID); err != nil {
			b.logger.Error("failed impersonate failed auth attempts", "error", err)
		}

		// send notification
		if _, err = b.notificationsService.FailLoginAttempts(user.UID); err != nil {
			b.logger.Error("Can't send notification", "error", err)
		}

		return true
	}

	return false
}

func (b Blocker) loginSecuritySettings() *syssettings.LoginSecuritySettings {
	return b.loginSettings
}

func (b *Blocker) createFailAttempt(ip, uid string) {
	err := b.failAuthAttemptRepository.Create(ip, uid)
	if err != nil {
		b.logger.Error("failed to create failAuthAttempt", "error", err)
	}
}

func (b *Blocker) deleteOldAttempts(uid string) {
	timeUntil := b.createUntilTime(-userBlockingDurationMin)
	b.failAuthAttemptRepository.DeleteAllOld(uid, timeUntil)
}

func (b *Blocker) createUntilTime(durationAsMinutes int64) time.Time {
	return time.Now().Add(time.Minute * time.Duration(durationAsMinutes))
}
