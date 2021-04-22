package models

import (
	"time"
)

// AccessLog is the abstract accesslog model
type AccessLog struct {
	ALID      uint64    `gorm:"primary_key:yes;column:alid;unique_index" json:"alid"`
	UID       string    `gorm:"column:uid;not null;default:0;" json:"uid"`
	IP        string    `gorm:"column:ip;not null;default:0;" json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName sets AccessLog's table name to be `users_accesslog`
func (AccessLog) TableName() string {
	return "users_accesslog"
}

// New creates new model
func (AccessLog) New(uid string) *AccessLog {
	return &AccessLog{
		UID: uid,
		IP:  "",
	}
}
