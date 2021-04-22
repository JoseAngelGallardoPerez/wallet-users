package models

import "time"

const (
	AttributeTypeString = "string"
	AttributeTypeBool   = "bool"
	AttributeTypeInt    = "int"
	AttributeTypeFloat  = "float"
)

type UserAttributeValue struct {
	UserID      string      `gorm:"primary_key;auto_increment:false;column:user_id" json:"userId"`
	AttributeId uint64      `gorm:"primary_key;auto_increment:false;column:attribute_id" json:"attributeId"`
	Value       interface{} `gorm:"column:value" json:"value"`
	CreatedAt   time.Time   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time   `gorm:"column:updated_at" json:"updatedAt"`
}

type Attribute struct {
	Id          uint64    `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Slug        string    `gorm:"column:slug" json:"slug"`
	Type        string    `gorm:"column:type" json:"type"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
}
