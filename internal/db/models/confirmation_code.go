package models

import (
	"time"
)

const (
	ConfirmationCodeSubjectPasswordRecovery      = "password_recovery"
	ConfirmationCodeSubjectPhoneVerificationCode = "phone_verification"
	// ConfirmationCodeSubjectSetPassword is used to enable a user to set a password when admin creates a profile
	ConfirmationCodeSubjectSetPassword           = "set_password"
	ConfirmationCodeSubjectEmailVerificationCode = "email_verification"
)

type ConfirmationCode struct {
	ID        uint64    `gorm:"primary_key" json:"-"`
	Code      string    `json:"code"`
	UserUID   string    `json:"-"`
	User      *User     `gorm:"foreignkey:UserUID;association_foreignkey:UID;association_autoupdate:false" json:"-"`
	Subject   string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
}
