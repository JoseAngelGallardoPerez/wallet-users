package users

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
)

type AttributeService struct {
	attributeRepo          *repositories.AttributeRepository
	userAttributeValueRepo *repositories.UserAttributeValueRepository
}

func NewAttributeService(
	attributeRepo *repositories.AttributeRepository,
	userAttributeValueRepo *repositories.UserAttributeValueRepository,

) *AttributeService {
	return &AttributeService{
		attributeRepo,
		userAttributeValueRepo,
	}
}

// Attaches attributes to the user
func (s *AttributeService) AttachAttributes(attributes map[string]interface{}, userId string, tx *gorm.DB) error {
	if len(attributes) > 0 {
		attributeSlugs := make([]string, 0, len(attributes))
		for k := range attributes {
			attributeSlugs = append(attributeSlugs, k)
		}
		attributeModels, err := s.attributeRepo.FindBySlugs(attributeSlugs)
		if err != nil {
			return errors.Wrap(err, "cannot find attributes")
		}

		userAttributeValueRepo := s.userAttributeValueRepo.WrapContext(tx)

		for key, val := range attributes {
			attributeModel, err := s.findAttribute(attributeModels, key)
			if err != nil {
				return err
			}

			attr := &models.UserAttributeValue{
				UserID:      userId,
				AttributeId: attributeModel.Id,
				Value:       val,
			}

			if err := userAttributeValueRepo.Save(attr); err != nil {
				return err
			}
		}
	}

	return nil
}

// find attribute by slug in a slice
func (s *AttributeService) findAttribute(models []*models.Attribute, slug string) (*models.Attribute, error) {
	for _, model := range models {
		if model.Slug == slug {
			return model, nil
		}
	}

	return nil, errors.Errorf("cannot find attribute '%s' in DB", slug)
}
