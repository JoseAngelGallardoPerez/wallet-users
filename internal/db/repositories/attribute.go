package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type AttributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) *AttributeRepository {
	return &AttributeRepository{db}
}

func (r *AttributeRepository) FindBySlug(slug string) (*models.Attribute, error) {
	record := &models.Attribute{}
	err := r.db.Where("slug = ?", slug).First(&record).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return record, nil
}

func (r *AttributeRepository) FindBySlugs(attributes []string) ([]*models.Attribute, error) {
	var records []*models.Attribute
	err := r.db.Where("slug IN (?)", attributes).Find(&records).Error

	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r AttributeRepository) WrapContext(db *gorm.DB) *AttributeRepository {
	r.db = db
	return &r
}
