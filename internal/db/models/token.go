package models

import (
	"time"
)

type Token struct {
	ID             uint64 `gorm:"primary_key"`
	Subject        string
	SignedString   string
	UserUID        string
	User           *User `gorm:"foreignkey:UserUID;association_foreignkey:UID;association_autoupdate:false"`
	RefreshTokenId *uint64
	RefreshToken   *Token `gorm:"foreignkey:RefreshTokenId;association_foreignkey:ID;association_autoupdate:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
