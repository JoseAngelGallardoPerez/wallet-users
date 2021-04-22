package models

import "time"

type VerificationFile struct {
	ID             uint32    `json:"id"`
	VerificationID uint32    `json:"verificationId"`
	FileID         uint64    `json:"fileId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
