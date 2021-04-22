package userprofilesrow

import (
	"github.com/Confialink/wallet-users/internal/db/models"
)

type physicalAddressFields struct {
	physicalAddress *models.Address
}

func (f *physicalAddressFields) call() []string {
	return []string{
		f.address(),
		f.addressSecondLine(),
		f.city(),
		f.state(),
		f.zipCode(),
		f.country(),
	}
}

func (f *physicalAddressFields) address() string {
	return f.physicalAddress.Address
}

func (f *physicalAddressFields) addressSecondLine() string {
	return f.physicalAddress.AddressSecondLine
}

func (f *physicalAddressFields) city() string {
	return f.physicalAddress.City
}

func (f *physicalAddressFields) state() string {
	return f.physicalAddress.Region
}

func (f *physicalAddressFields) zipCode() string {
	return f.physicalAddress.ZipCode
}

func (f *physicalAddressFields) country() string {
	return f.physicalAddress.CountryIsoTwo
}
