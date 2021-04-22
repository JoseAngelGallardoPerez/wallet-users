package models

import "time"

// SecurityQuestionsAnswer contains possible security questions answers
type SecurityQuestionsAnswer struct {
	AID       uint64    `gorm:"primary_key:yes;column:aid;unique_index" json:"aid"`
	SQID      uint64    `gorm:"column:sqid;not null;" json:"sqid" binding:"required"`
	UID       string    `gorm:"column:uid;not null;default:0;" json:"-"`
	Answer    string    `gorm:"column:answer" json:"answer" binding:"required,max=255"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
