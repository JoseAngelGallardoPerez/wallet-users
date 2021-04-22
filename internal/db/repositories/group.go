package repositories

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

// UserGroupsRepositoryHandler is interface for repository functionality that ought to be implemented manually.
type UserGroupsRepositoryHandler interface {
	FindById(id uint64) (*models.UserGroup, error)
	FindByName(name string) (*models.UserGroup, error)
	Create(userGroup *models.UserGroup) error
	Update(userGroup *models.UserGroup, data *models.UserGroup) error
	Delete(userGroup *models.UserGroup) error
	Filter(params url.Values) *gorm.DB
	Paginate(query *gorm.DB, pageQuery string, limitQuery string) (*Pagination, error)
}

// UserGroupsRepository is userGroup group repository for CRUD operations.
type UserGroupsRepository struct {
	DB *gorm.DB
}

func NewUserGroupsRepository(db *gorm.DB) *UserGroupsRepository {
	return &UserGroupsRepository{
		db,
	}
}

// FindById find user group by id
func (repo *UserGroupsRepository) FindById(id uint64) (*models.UserGroup, error) {
	userGroup := &models.UserGroup{}
	if err := repo.DB.Where("id = ?", id).
		First(&userGroup).Error; err != nil {
		return nil, err
	}
	return userGroup, nil
}

// FindByName find user group by name
func (repo *UserGroupsRepository) FindByName(name string) (*models.UserGroup, error) {
	userGroup := &models.UserGroup{}
	if err := repo.DB.Where("name = ?", name).
		First(&userGroup).Error; err != nil {
		return nil, err
	}
	return userGroup, nil
}

// Create creates new user
func (repo *UserGroupsRepository) Create(userGroup *models.UserGroup) error {
	if err := repo.DB.Create(userGroup).Error; err != nil {
		return err
	}
	return nil
}

// Update updates an existing user
func (repo *UserGroupsRepository) Update(userGroup *models.UserGroup, data *models.UserGroup) error {
	if err := repo.DB.Model(&userGroup).Save(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete delete an existing user
func (repo *UserGroupsRepository) Delete(userGroup *models.UserGroup) error {
	if err := repo.DB.Delete(userGroup).Error; err != nil {
		return err
	}
	return nil
}

// Filter apply request params to the builder instance.
func (repo *UserGroupsRepository) Filter(params url.Values) *gorm.DB {
	db := repo.DB
	db = repo.applyFilters(db, params)
	db = repo.applySort(db, params)
	return db
}

// Paginate returns a new Pagination instance.
func (repo *UserGroupsRepository) Paginate(query *gorm.DB, pageQuery string, limitQuery string) (*Pagination, error) {
	p := &Pagination{}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return p, errors.New("Invalid parameter")
	}
	p.Limit = int(math.Max(1, math.Min(10000, float64(limit))))

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return p, errors.New("Invalid parameter")
	}
	p.Page = int(math.Max(1, float64(page)))

	p.Offset = p.Limit * (p.Page - 1)

	done := make(chan bool, 1)

	var userGroups []*models.UserGroup
	var count int

	go countItems(query, userGroups, done, &count)

	if err := query.Limit(p.Limit).Offset(p.Offset).Find(&userGroups).Error; err != nil {
		return nil, err
	}
	<-done

	p.TotalRecord = count
	p.Items = userGroups
	p.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	return p, nil
}

func (repo *UserGroupsRepository) applyFilters(query *gorm.DB, params url.Values) *gorm.DB {
	if len(params.Get("filter[query]")) > 0 {
		value := "%" + params.Get("filter[query]") + "%"
		query = query.Where("name LIKE ?", value)
	}
	return query
}

func (repo *UserGroupsRepository) applySort(query *gorm.DB, params url.Values) *gorm.DB {
	order := "created_at desc"
	if len(params.Get("sort")) > 0 {
		field := params.Get("sort")
		if (string)(field[0]) == "-" {
			order = fmt.Sprintf("%s %s", field[1:], "desc")
		} else {
			order = fmt.Sprintf("%s %s", field, "asc")
		}
	}
	query = query.Order(order)
	return query
}
