package repositories

import (
	"fmt"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

// InvitesRepository is invites repository for CRUD operations.
type InvitesRepository struct {
	DB *gorm.DB
}

func NewInvitesRepository(db *gorm.DB) *InvitesRepository {
	return &InvitesRepository{
		db,
	}
}

// FindByCode find invite by code
func (repo *InvitesRepository) FindByCode(name string) (*models.Invite, error) {
	invite := &models.Invite{}
	if err := repo.DB.Where("code = ?", name).
		First(&invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

// FindByUserUID find user by user UID
func (repo *InvitesRepository) FindByUserUID(userUID string) ([]*models.Invite, error) {
	var invites []*models.Invite
	if err := repo.DB.Where("user_uid = ?", userUID).Find(&invites).Error; err != nil {
		return nil, fmt.Errorf("could not find invites by user uid `%s` in database", userUID)
	}
	return invites, nil
}

// CountByUserUID count invites by user uid
func (repo *InvitesRepository) CountByUserUID(userUID string) (uint64, error) {
	var invite models.Invite
	var count uint64
	if err := repo.DB.Where("user_uid = ?", userUID).Model(&invite).Count(&count).
		Error; err != nil {
		return count, err
	}
	return count, nil
}

// Create creates a new invite
func (r *InvitesRepository) Create(invite *models.Invite) error {
	if err := r.DB.Create(invite).Error; err != nil {
		return err
	}
	return nil
}

// Update updates an existing invite
func (r *InvitesRepository) Update(invite *models.Invite, data *models.Invite) error {
	if err := r.DB.Model(&invite).Save(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete delete an existing invite
func (r *InvitesRepository) Delete(invite *models.Invite) error {
	if err := r.DB.Delete(invite).Error; err != nil {
		return err
	}
	return nil
}

func (copy InvitesRepository) WrapContext(db *gorm.DB) *InvitesRepository {
	copy.DB = db
	return &copy
}
