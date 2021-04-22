package users

import (
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/serializers"
)

type UserLoaderService struct {
	addressRepo            *repositories.AddressRepository
	attributeRepo          *repositories.AttributeRepository
	userAttributeValueRepo *repositories.UserAttributeValueRepository
}

func NewUserLoaderService(
	addressRepo *repositories.AddressRepository,
	attributeRepo *repositories.AttributeRepository,
	userAttributeValueRepo *repositories.UserAttributeValueRepository,
) *UserLoaderService {
	return &UserLoaderService{
		addressRepo,
		attributeRepo,
		userAttributeValueRepo,
	}
}

// LoadUserCompletely loads additional user data into the user structure
//
// PhysicalAddresses ( from `addresses` table )
// MailingAddresses ( from `addresses` table )
// Attributes ( from `attributes` table )
func (u *UserLoaderService) LoadUserCompletely(user *models.User) error {
	mailingAddresses, err := u.addressRepo.FindMailingByUser(user.UID)
	if err != nil {
		return errors.Wrap(err, "UserLoaderService: сan't load mailing addresses")
	}

	user.MailingAddresses = mailingAddresses

	physicalAddresses, err := u.addressRepo.FindPhysicalByUser(user.UID)
	if err != nil {
		return errors.Wrap(err, "UserLoaderService: сan't load physical addresses")
	}

	user.PhysicalAddresses = physicalAddresses

	rawAttributes, err := u.userAttributeValueRepo.AllByUserId(user.UID)
	if err != nil {
		return errors.Wrap(err, "UserLoaderService: can't load attributes")
	}

	userAttributes := make(map[string]interface{})

	for _, rawAttribute := range rawAttributes {
		userAttributes[rawAttribute.Slug] = ToTypedValue(rawAttribute.Value, rawAttribute.Type)
	}

	user.Attributes = userAttributes

	return nil
}

// LoadUserCompletelyAndSerialize loads additional user data into the user structure,
// and than serialize the user structure to JSON-ready output format
func (u *UserLoaderService) LoadUserCompletelyAndSerialize(user *models.User) (map[string]interface{}, error) {
	err := u.LoadUserCompletely(user)
	if err != nil {
		return nil, err
	}

	return serializers.NewGetUser(user).Serialize(), nil
}
