package models

type Role struct {
	Slug string `gorm:"column:slug" json:"slug"`
	Name string `gorm:"column:name" json:"name"`
}
