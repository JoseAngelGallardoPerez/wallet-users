package repositories

import (
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

type FailAuthAttemptRepository struct {
	DB *gorm.DB
}

func NewFailAuthAttemptRepository(db *gorm.DB) *FailAuthAttemptRepository {
	return &FailAuthAttemptRepository{
		db,
	}
}

func (repo *FailAuthAttemptRepository) Create(ip string, uid string) error {
	if err := repo.DB.Exec("INSERT INTO fail_auth_attempts (`ip`, `uid`, `created_at`) VALUES (INET_ATON(?), ?, CURRENT_TIMESTAMP())", ip, uid).Error; err != nil {
		return err
	}

	return nil
}

func (repo *FailAuthAttemptRepository) GetCountByIp(ip string, timeFrom time.Time) uint32 {
	var count uint32
	repo.DB.Model(&models.FailAuthAttempt{}).Where("ip = INET_ATON(?) AND created_at >= ?", ip, timeFrom).Count(&count)

	return count
}

func (repo *FailAuthAttemptRepository) GetCountByUID(uid string, timeFrom time.Time) uint32 {
	var count uint32
	repo.DB.Model(&models.FailAuthAttempt{}).Where("uid = ? AND created_at >= ?", uid, timeFrom).Count(&count)

	return count
}

func (repo *FailAuthAttemptRepository) DeleteAllByIp(ip string) error {
	if err := repo.DB.Where("ip = INET_ATON(?)", ip).Delete(models.FailAuthAttempt{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *FailAuthAttemptRepository) DeleteAllOld(uid string, time time.Time) {
	repo.DB.Where("uid = ? OR created_at <= ?", uid, time).Delete(models.FailAuthAttempt{})
}

func (repo *FailAuthAttemptRepository) ImpersonateAllByUID(uid string) error {
	if err := repo.DB.Exec("UPDATE fail_auth_attempts SET `uid` = '' WHERE `uid` = ?", uid).Error; err != nil {
		return err
	}

	return nil
}
