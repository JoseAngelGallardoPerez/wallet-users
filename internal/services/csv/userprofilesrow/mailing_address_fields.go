package userprofilesrow

import (
	"github.com/Confialink/wallet-users/internal/db/models"
)

type mailingAddressFields struct {
	mailingAddress *models.Address
}

func (f *mailingAddressFields) call() []string {
	return []string{
		f.sameAsPhysical(),
		f.name(),
		f.address(),
		f.addressSecondLine(),
		f.city(),
		f.state(),
		f.zipCode(),
		f.country(),
		f.phoneNumber(),
	}
}

func (f *mailingAddressFields) sameAsPhysical() string {
	return "False" // todo: refactor
}

func (f *mailingAddressFields) name() string {
	return f.mailingAddress.Name
}

func (f *mailingAddressFields) address() string {
	return f.mailingAddress.Address
}

func (f *mailingAddressFields) addressSecondLine() string {
	return f.mailingAddress.AddressSecondLine
}

func (f *mailingAddressFields) city() string {
	return f.mailingAddress.City
}

func (f *mailingAddressFields) state() string {
	return f.mailingAddress.Region
}

func (f *mailingAddressFields) zipCode() string {
	return f.mailingAddress.ZipCode
}

func (f *mailingAddressFields) country() string {
	return f.mailingAddress.CountryIsoTwo
}

func (f *mailingAddressFields) phoneNumber() string {
	return f.mailingAddress.PhoneNumber
}
