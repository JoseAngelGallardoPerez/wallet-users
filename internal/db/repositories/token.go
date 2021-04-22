package repositories

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

// UserGroupsRepository is userGroup group repository for CRUD operations.
type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(DB *gorm.DB) *TokenRepository {
	return &TokenRepository{DB: DB}
}

func (repo *TokenRepository) FindTokenBySignedString(signedString string) (*models.Token, error) {
	model := &models.Token{}
	if err := repo.DB.Where("signed_string = ?", signedString).
		Preload("User").
		First(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (repo *TokenRepository) FindTokenBySignedStringAndSubject(signedString string, subject string) (*models.Token, error) {
	model := &models.Token{}
	if err := repo.DB.Where("signed_string = ? AND subject = ?", signedString, subject).
		Preload("User").
		First(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (repo *TokenRepository) FindTokensBySubject(subject string) ([]*models.Token, error) {
	var tokens []*models.Token
	if err := repo.DB.Where("subject = ?", subject).
		Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (repo *TokenRepository) FindAccessTokenByRefreshTokenId(refreshTokenId uint64) (*models.Token, error) {
	var token *models.Token
	if err := repo.DB.Where("refresh_token_id = ?", refreshTokenId).
		First(&token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (repo *TokenRepository) DeleteTokensByUID(uid string) error {
	if err := repo.DB.Where("user_uid = ?", uid).
		Delete(&models.Token{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo *TokenRepository) DeleteTokenByID(id uint64) error {
	if err := repo.DB.Where("id = ?", id).
		Delete(&models.Token{}).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo *TokenRepository) Create(token *models.Token) (*models.Token, error) {
	if err := repo.DB.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (repo *TokenRepository) Delete(token *models.Token) error {
	if err := repo.DB.Delete(token).Error; err != nil {
		return err
	}
	return nil
}
