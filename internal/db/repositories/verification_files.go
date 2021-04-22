package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type VerificationFilesRepository struct {
	db *gorm.DB
}

func NewVerificationFilesRepository(db *gorm.DB) *VerificationFilesRepository {
	return &VerificationFilesRepository{db}
}

func (r *VerificationFilesRepository) Create(file *models.VerificationFile) (*models.VerificationFile, error) {
	return file, r.db.Create(file).Error
}

func (r *VerificationFilesRepository) FindByVerificationID(id uint32) ([]*models.VerificationFile, error) {
	var files []*models.VerificationFile
	return files, r.db.Where("verification_id = ?", id).Find(&files).Error
}

func (r *VerificationFilesRepository) CountByVerificationID(id uint32) (int, error) {
	var count int
	return count, r.db.Model(models.VerificationFile{}).Where("verification_id = ?", id).Count(&count).Error
}

func (r VerificationFilesRepository) WrapContext(db *gorm.DB) *VerificationFilesRepository {
	r.db = db
	return &r
}
