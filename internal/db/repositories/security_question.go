package repositories

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

// SecurityQuestionRepositoryHandler is interface for repository functionality that ought to be implemented manually.
type SecurityQuestionRepositoryHandler interface {
	GetAll() (*models.SecurityQuestion, error)
	Create(question *models.SecurityQuestion) error
	Update(question *models.SecurityQuestion, data *models.SecurityQuestion) error
	Delete(question *models.SecurityQuestion) error
}

// SecurityQuestionRepository is security question repository for CRUD operations.
type SecurityQuestionRepository struct {
	DB *gorm.DB
}

func NewSecurityQuestionRepository(db *gorm.DB) *SecurityQuestionRepository {
	return &SecurityQuestionRepository{
		db,
	}
}

// GetAll returns all security questions
func (repo *SecurityQuestionRepository) GetAll() ([]*models.SecurityQuestion, error) {
	var questions []*models.SecurityQuestion
	if err := repo.DB.Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// FindByUID returns list of security questions by uid
func (repo *SecurityQuestionRepository) FindByUID(uid string) ([]*models.SecurityQuestion, error) {
	var questions []*models.SecurityQuestion
	if err := repo.DB.Where("uid = ?", uid).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// FindBySqids returns list of security questions by its ids
func (repo *SecurityQuestionRepository) FindBySqids(sqids []uint64) ([]*models.SecurityQuestion, error) {
	var questions []*models.SecurityQuestion
	if err := repo.DB.Where("sqid IN (?)", sqids).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

// Create creates new security question
func (repo *SecurityQuestionRepository) Create(question *models.SecurityQuestion) error {
	if err := repo.DB.Debug().Create(question).Error; err != nil {
		return err
	}
	return nil
}

// Update updates an existing security question
func (repo *SecurityQuestionRepository) Update(question *models.SecurityQuestion, data *models.SecurityQuestion) error {
	if err := repo.DB.Model(&question).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete delete an existing security question
func (repo *SecurityQuestionRepository) Delete(question *models.SecurityQuestion) error {
	if err := repo.DB.Delete(question).Error; err != nil {
		return err
	}
	return nil
}
