package repositories

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

// SecurityQuestionsAnswerRepository is answer repository for CRUD operations.
type SecurityQuestionsAnswerRepository struct {
	DB *gorm.DB
}

func NewSecurityQuestionsAnswerRepository(db *gorm.DB) *SecurityQuestionsAnswerRepository {
	return &SecurityQuestionsAnswerRepository{
		db,
	}
}

// FindByUID returns list of settings
func (r *SecurityQuestionsAnswerRepository) FindByUID(uid string) ([]*models.SecurityQuestionsAnswer, error) {
	var answers []*models.SecurityQuestionsAnswer
	query := r.DB.Where("uid = ?", uid)
	// user can have only one question and answer
	query = query.Order("updated_at DESC").Limit(1)
	if err := query.Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *SecurityQuestionsAnswerRepository) DeleteByUID(uid string) {
	r.DB.Where("uid = ?", uid).Delete(models.SecurityQuestionsAnswer{})
}

// FirstOrCreate updates or create if answer don't exists
func (r *SecurityQuestionsAnswerRepository) FirstOrCreate(data *models.SecurityQuestionsAnswer) (*models.SecurityQuestionsAnswer, error) {
	err := r.DB.Where(models.SecurityQuestionsAnswer{SQID: data.SQID, UID: data.UID}).
		Assign(models.SecurityQuestionsAnswer{
			UID:    data.UID,
			SQID:   data.SQID,
			Answer: data.Answer,
		}).
		FirstOrCreate(&data).Error
	if nil != err {
		return nil, err
	}
	return data, nil
}

func (copy SecurityQuestionsAnswerRepository) WrapContext(db *gorm.DB) *SecurityQuestionsAnswerRepository {
	copy.DB = db
	return &copy
}
