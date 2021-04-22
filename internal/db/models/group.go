package models

import (
	"time"
)

// UserGroup is the abstract user group model
// Models should only be concerned with database schema, more strict checking should be put in validator.
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
// NOTE: If you want to split null and "", you should use *string instead of string.
type UserGroup struct {
	ID          uint64    `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"unique_index" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
