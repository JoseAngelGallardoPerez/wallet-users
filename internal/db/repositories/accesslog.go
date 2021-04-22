package repositories

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

// AccesLogRepositoryHandler is interface for repository functionality that ought to be implemented manually.
type AccesLogRepositoryHandler interface {
	FindByUID(uid string) (*models.AccessLog, error)
	Create(accessLog *models.AccessLog) error
}

// AccesLogRepository is userGroup group repository for CRUD operations.
type AccesLogRepository struct {
	DB *gorm.DB
}

func NewAccesLogRepository(DB *gorm.DB) *AccesLogRepository {
	return &AccesLogRepository{DB: DB}
}

// FindByUID find records by uid
func (repo *AccesLogRepository) FindByUID(uid string) (*models.AccessLog, error) {
	accessLog := &models.AccessLog{}
	if err := repo.DB.Raw("SELECT `alid`, `uid`, INET_NTOA(`ip`) FROM users_accesslog WHERE `uid` = ?", uid).
		Scan(&accessLog).Error; err != nil {
		return nil, err
	}
	return accessLog, nil
}

// Create creates new record
func (repo *AccesLogRepository) Create(accessLog *models.AccessLog) error {
	if err := repo.DB.Exec("INSERT INTO users_accesslog (`uid`, `ip`, `created_at`) VALUES (?, INET_ATON(?), CURRENT_TIMESTAMP())", accessLog.UID, accessLog.IP).
		Error; err != nil {
		return err
	}
	return nil
}
