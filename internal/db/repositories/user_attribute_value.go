package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type UserAttributeValueRepository struct {
	db *gorm.DB
}

func NewUserAttributeValueRepository(db *gorm.DB) *UserAttributeValueRepository {
	return &UserAttributeValueRepository{db}
}

func (r *UserAttributeValueRepository) Save(value *models.UserAttributeValue) error {
	if err := r.db.Save(value).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserAttributeValueRepository) Delete(userId string, attributeId uint64) error {
	if err := r.db.Where("user_id = ? AND attribute_id = ?", userId, attributeId).
		Delete(&models.UserAttributeValue{}).Error; err != nil {
		return err
	}
	return nil
}

type RawUserAttributes struct {
	Slug  string
	Value string
	Type  string
}

func (r *UserAttributeValueRepository) AllByUserId(userId string) ([]*RawUserAttributes, error) {
	var attributes []*RawUserAttributes

	queryString := "SELECT a.slug as slug, uav.value as value, a.type as type  " +
		"FROM user_attribute_values as uav " +
		"INNER JOIN `attributes` as a on a.id = uav.attribute_id " +
		"WHERE uav.user_id = ?"

	err := r.db.Raw(queryString, userId).Scan(&attributes).Error

	if err != nil {
		return nil, err
	}

	return attributes, nil
}

func (r UserAttributeValueRepository) WrapContext(db *gorm.DB) *UserAttributeValueRepository {
	r.db = db
	return &r
}
