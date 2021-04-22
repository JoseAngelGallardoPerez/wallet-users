package models

import (
	"time"
)

type FailAuthAttempt struct {
	ID        uint64    `gorm:"primary_key:yes;column:id" json:"id"`
	IP        string    `gorm:"column:ip;not null;default:0;" json:"ip"`
	UID       string    `gorm:"column:uid" json:"uid"`
	CreatedAt time.Time `json:"createdAt"`
}

func (*FailAuthAttempt) TableName() string {
	return "fail_auth_attempts"
}
