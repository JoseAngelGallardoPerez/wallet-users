package models

import (
	"time"
)

const (
	VerificationStatusPending   = "pending"
	VerificationStatusProgress  = "progress"
	VerificationStatusApproved  = "approved"
	VerificationStatusCancelled = "cancelled"
)

const (
	VerificationTypePersonalId    = "personal_id"
	VerificationTypePersonalPhoto = "personal_photo"
)

type Verification struct {
	ID        uint32              `gorm:"primary_key" json:"id"`
	Status    string              `gorm:"column:status; default:pending" json:"status" json.enum:"pending,progress,approved,cancelled"`
	Type      string              `gorm:"column:type; default:personal_id" json:"type" json.enum:"personal_id,personal_photo"`
	Files     []*VerificationFile `json:"files"`
	UserUID   string              `json:"-"`
	User      *User               `gorm:"foreignkey:UserUID;association_foreignkey:UID;association_autoupdate:false" json:"-"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

func GetAvailableVerificationTypes() map[string]bool {
	return map[string]bool{
		VerificationTypePersonalId:    true,
		VerificationTypePersonalPhoto: true,
	}
}
