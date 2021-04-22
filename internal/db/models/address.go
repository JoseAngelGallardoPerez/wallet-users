package models

import "time"

const (
	AddressTypeMailing  = "mailing"
	AddressTypePhysical = "physical"
)

type Address struct {
	ID                uint64 `gorm:"column:id" json:"id"`
	UserID            string `gorm:"column:user_id" json:"userId"`
	Type              string `gorm:"column:type" json:"type"`
	CountryIsoTwo     string `gorm:"column:country_iso_two" json:"countryIsoTwo"`
	Region            string `gorm:"column:region" json:"region"`
	City              string `gorm:"column:city" json:"city"`
	ZipCode           string `gorm:"column:zip_code" json:"zipCode"`
	Address           string `gorm:"column:address" json:"address"`
	AddressSecondLine string `gorm:"column:address_second_line" json:"addressSecondLine"`
	Name              string `gorm:"column:name" json:"name"`
	PhoneNumber       string `gorm:"column:phone_number" json:"phoneNumber"`
	Description       string `gorm:"column:description" json:"description"`

	// Map
	Latitude  *float64 `gorm:"column:latitude" json:"latitude"`
	Longitude *float64 `gorm:"column:longitude" json:"longitude"`

	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
