package models

// SecurityQuestion contains possible security questions
// NOTE: UID maybe equal 0 for questions available system-wide, or the owning uid for custom per-user questions
type SecurityQuestion struct {
	SQID     uint64 `gorm:"primary_key:yes;column:sqid;unique_index" json:"sqid"`
	UID      string `gorm:"column:uid;not null;default:0;" json:"uid"`
	Question string `gorm:"column:question" json:"question"`
}
