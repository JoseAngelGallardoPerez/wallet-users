package models

import (
	"time"
)

// BlockedIp is the abstract blocked ips model
type BlockedIp struct {
	ID           uint64    `gorm:"primary_key:yes;column:id;unique_index" json:"id"`
	IP           string    `gorm:"column:ip;not null;default:0;" json:"ip"`
	CreatedAt    time.Time `json:"createdAt"`
	BlockedUntil time.Time `json:"blockedUntil"`
}

// TableName sets BlockedIp's table name to be `blocked_ips`
func (BlockedIp) TableName() string {
	return "blocked_ips"
}
