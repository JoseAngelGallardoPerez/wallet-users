package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type FormRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) *FormRepository {
	return &FormRepository{db}
}

// Returns all forms
func (repo *FormRepository) All() ([]*models.Form, error) {
	var forms []*models.Form
	if err := repo.db.Find(&forms).Error; err != nil {
		return nil, err
	}
	return forms, nil
}
