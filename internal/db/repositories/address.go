package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db}
}

func (r *AddressRepository) Save(address *models.Address) error {
	if err := r.db.Save(address).Error; err != nil {
		return err
	}
	return nil
}

func (r *AddressRepository) FindMailingByUser(userId string) ([]*models.Address, error) {
	return r.findByUserAndType(userId, models.AddressTypeMailing)
}

func (r *AddressRepository) FindPhysicalByUser(userId string) ([]*models.Address, error) {
	return r.findByUserAndType(userId, models.AddressTypePhysical)
}

func (r *AddressRepository) findByUserAndType(userId, addressType string) ([]*models.Address, error) {
	var records []*models.Address
	if err := r.db.Where("user_id = ? AND type = ?", userId, addressType).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("could not find addresses with user_id `%s` in database", userId)
	}
	return records, nil
}

func (r *AddressRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Address{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *AddressRepository) DeleteByTypeForUser(userId, addrType string) error {
	if err := r.db.Where("user_id = ? AND type = ?", userId, addrType).Delete(&models.Address{}).Error; err != nil {
		return err
	}
	return nil
}

func (r AddressRepository) WrapContext(db *gorm.DB) *AddressRepository {
	r.db = db
	return &r
}
