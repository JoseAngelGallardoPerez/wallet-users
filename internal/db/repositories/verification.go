package repositories

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

type VerificationRepository struct {
	db *gorm.DB
}

// NewVerificationRepository creates new repository
func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db}
}

func (r *VerificationRepository) FindById(id uint64) (*models.Verification, error) {
	verification := &models.Verification{}
	return verification, r.db.Preload("User").Preload("Files").Find(verification, id).Error
}

func (r *VerificationRepository) FindByUID(uid string) ([]*models.Verification, error) {
	var verifications []*models.Verification
	if err := r.db.Where("user_uid = ?", uid).
		Preload("Files").
		Find(&verifications).
		Error; err != nil {
		return nil, err
	}
	return verifications, nil
}

func (r *VerificationRepository) FindByUIDAndType(uid, verificationType string) (*models.Verification, error) {
	verification := models.Verification{}
	err := r.db.Where("user_uid = ? AND type = ?", uid, verificationType).
		First(&verification).
		Error
	return &verification, err
}

func (r *VerificationRepository) Create(verification *models.Verification) (*models.Verification, error) {
	if err := r.db.Create(verification).Error; err != nil {
		return nil, err
	}
	return verification, nil
}

func (r *VerificationRepository) Save(verification *models.Verification) error {
	if err := r.db.Save(verification).Error; err != nil {
		return err
	}
	return nil
}

func (r VerificationRepository) WrapContext(db *gorm.DB) *VerificationRepository {
	r.db = db
	return &r
}
