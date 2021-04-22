package repositories

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/Confialink/wallet-pkg-list_params/adapters"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/http/serializers"
)

// UsersRepository is user repository for CRUD operations.
type UsersRepository struct {
	DB *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{
		db,
	}
}

// Paginate returns a new Pagination instance.
func (repo *UsersRepository) Paginate(
	query *gorm.DB,
	pageQuery string,
	limitQuery string,
	serializer serializers.UsersCollectionSerializer,
) (*Pagination, error) {
	p := &Pagination{}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return p, errors.New("invalid parameter")
	}
	p.Limit = int(math.Max(1, math.Min(10000, float64(limit))))

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return p, errors.New("invalid parameter")
	}
	p.Page = int(math.Max(1, float64(page)))

	p.Offset = p.Limit * (p.Page - 1)

	done := make(chan bool, 1)

	var users []*models.User
	var count int

	go countItems(query, users, done, &count)

	if err := query.Limit(p.Limit).Offset(p.Offset).Find(&users).Error; err != nil {
		return nil, err
	}
	<-done

	p.TotalRecord = count
	p.Items = serializer.Serialize(users)
	p.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	return p, nil
}

// Filter apply request params to the builder instance.
func (repo *UsersRepository) Filter(params url.Values) *gorm.DB {
	db := repo.DB
	db = db.Preload("UserGroup")
	db = db.Preload("CompanyDetails")
	db = repo.applyFilters(db, params)
	db = repo.applySort(db, params)
	return db
}

func (repo *UsersRepository) FilterByCompanyName(query *gorm.DB, companyName string) *gorm.DB {
	query = query.Where("company_name = ?", companyName)
	return query
}

func (repo *UsersRepository) FilterByRoleName(query *gorm.DB, roleName string) *gorm.DB {
	query = query.Where("role_name = ?", roleName)
	return query
}

func (repo *UsersRepository) applyFilters(query *gorm.DB, params url.Values) *gorm.DB {

	// Filter by query
	if len(params.Get("filter[query]")) > 0 {
		columns := []string{"uid", "email", "username", "first_name", "last_name"}
		conditions, values := repo.processSearchQueryParam(params.Get("filter[query]"), columns)
		if len(conditions) > 0 {
			query = query.Where(strings.Join(conditions, " OR "), values...)
		}
	}

	// Filter by status
	if len(params.Get("filter[status]")) > 0 {
		query = query.Where("status = ?", params.Get("filter[status]"))
	}

	// Filter by blockedUntil
	if len(params.Get("filter[isBlocked]")) > 0 {
		if params.Get("filter[isBlocked]") == "true" {
			query = query.Where("blocked_until > ?", time.Now())
		} else {
			query = query.Where("blocked_until <= ? OR blocked_until IS NULL", time.Now())
		}
	}

	// Filter by role
	roles := params["filter[role_name]"]
	if len(roles) > 0 {
		query = query.Where("role_name IN (?)", roles)
	}

	// Filter by group
	if _, err := strconv.ParseInt(params.Get("filter[user_group_id]"), 10, 64); err == nil {
		query = query.Where("user_group_id = ?", params.Get("filter[user_group_id]"))
	}

	// Filter by dates
	if len(params.Get("filter[date_from]")) > 0 {
		query = query.Where("created_at > ?", params.Get("filter[date_from]"))
	}

	if len(params.Get("filter[date_to]")) > 0 {
		query = query.Where("created_at < ?", params.Get("filter[date_to]"))
	}

	return query
}

func (repo *UsersRepository) applySort(query *gorm.DB, params url.Values) *gorm.DB {
	order := "created_at desc"
	if len(params.Get("sort")) > 0 {
		field := params.Get("sort")
		if strings.Contains(field, "companies") {
			query = query.Joins("JOIN companies ON companies.id = users.company_id")
		}
		if (string)(field[0]) == "-" {
			order = fmt.Sprintf("%s %s", field[1:], "desc")
		} else {
			order = fmt.Sprintf("%s %s", field, "asc")
		}
	}
	query = query.Order(order)
	return query
}

// FindByUID find user by uid
func (repo *UsersRepository) FindByUID(uid string) (*models.User, error) {
	user := &models.User{}
	if err := repo.DB.Where("uid = ?", uid).Preload("UserGroup").Preload("CompanyDetails").
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("could not find user with uid `%s` in database", uid)
	}
	return user, nil
}

// FindByUIDWithoutGroup find user by uid without preloaded group
func (repo *UsersRepository) FindByUIDWithoutGroup(uid string) (*models.User, error) {
	user := &models.User{}
	if err := repo.DB.Where("uid = ?", uid).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("could not find user with uid `%s` in database", uid)
	}
	return user, nil
}

// FindByUID find user by uid
func (repo *UsersRepository) FindByClassId(classId uint64) ([]*models.User, error) {
	var users []*models.User
	if err := repo.DB.Where("class_id = ?", classId).Preload("UserGroup").Preload("CompanyDetails").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("could not find user with class_id `%d` in database", classId)
	}
	return users, nil
}

// FindByUsername find user by username
func (repo *UsersRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	// TODO: move preloads to a parameter
	if err := repo.DB.Where("username = ?", username).Preload("UserGroup").Preload("CompanyDetails").
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByPhoneNumber find user by phone_number
func (repo *UsersRepository) FindByPhoneNumber(phoneNumber string) (*models.User, error) {
	user := &models.User{}
	if err := repo.DB.Where("phone_number = ?", phoneNumber).Preload("UserGroup").
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindActiveByPhoneNumbers find active users by phone_number
func (repo *UsersRepository) FindActiveByPhoneNumbers(phoneNumber []string) ([]*models.User, error) {
	var users []*models.User

	if err := repo.DB.Where("phone_number IN (?)", phoneNumber).
		Where("status = ?", models.StatusActive).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("could not find user with phone_number `%s` in database", phoneNumber)
	}
	return users, nil
}

// FindByProfileData find users
func (repo *UsersRepository) FindByProfileData(searchQuery string, fields []string) ([]*models.User, error) {
	conditions, values := repo.processSearchQueryParam(searchQuery, fields)
	query := repo.DB.Where(strings.Join(conditions, " OR "), values...)

	var users []*models.User
	if err := query.Limit(100).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FindByEmail find user by email
func (repo *UsersRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := repo.DB.Where("email = ?", email).Preload("UserGroup").Preload("CompanyDetails").
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByEmail find user by email
func (repo *UsersRepository) IsExistByEmailWithoutCurrent(email, currentUserUID string) (bool, error) {
	var count uint64
	var user models.User
	if err := repo.DB.Where("email = ?", email).
		Where("uid != ?", currentUserUID).
		Model(&user).Count(&count).Error; err != nil {
		return false, fmt.Errorf("could not get users count with email `%s` in database. Error: %s", email, err)
	}
	return count > 0, nil
}

// FindByEmailOrPhoneNumber find user by email or username
func (repo *UsersRepository) FindByEmailOrPhoneNumber(value string) (*models.User, error) {
	user := &models.User{}
	if err := repo.DB.Where("email = ? OR phone_number = ?", value, value).Preload("UserGroup").Preload("CompanyDetails").
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("could not find user with email or username `%s` in database", value)
	}
	return user, nil
}

// FindByRoleName find user by rolename
func (repo *UsersRepository) FindByRoleName(rolename string) ([]*models.User, error) {
	var users []*models.User
	db := repo.DB.Where("role_name = ?", rolename)

	if err := db.Preload("CompanyDetails").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindUserByTokenAndSubject find user by token and subject
func (repo *UsersRepository) FindUserByTokenAndSubject(token string, subject string) (*models.User, error) {
	user := &models.User{}
	// TODO: move preloads to a parameter
	query := repo.DB.Joins("LEFT JOIN tokens ON tokens.user_uid = users.uid").
		Where("tokens.signed_string = ? AND tokens.subject = ?", token, subject).
		Preload("UserGroup").
		Preload("CompanyDetails").
		First(&user)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("can not find user by access token `%s` in database", subject)
	}
	return user, nil
}

// GetByUIDs returns users by passed uids
func (repo *UsersRepository) GetByUIDs(uids []string) ([]*models.User, error) {
	var users []*models.User
	if err := repo.DB.Where("uid IN (?)", uids).Preload("CompanyDetails").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByParentUID returns users by parent uid
func (repo *UsersRepository) GetByParentUID(parentUID string) ([]*models.User, error) {
	var users []*models.User
	if err := repo.DB.Where("parent_id = ?", parentUID).Preload("CompanyDetails").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByUserGroupId returns users by passed group id
func (repo *UsersRepository) GetByUserGroupId(id uint64) ([]*models.User, error) {
	var users []*models.User
	if err := repo.DB.Where("user_group_id = ?", id).Preload("CompanyDetails").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetAll returns all users
func (repo *UsersRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	if err := repo.DB.Preload("CompanyDetails").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Create creates new user
func (repo *UsersRepository) Create(user *models.User, confirmed bool) (*models.User, error) {
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return repo.FindByUID(user.UID)
}

// Update updates an existing user
func (repo *UsersRepository) Update(user *models.User) (*models.User, error) {
	if err := repo.DB.Model(&models.User{UID: user.UID}).Updates(user).Error; err != nil {
		return nil, err
	}
	return repo.FindByUID(user.UID)
}

// Save all fields an existing user
func (repo *UsersRepository) Save(user *models.User) (*models.User, error) {
	if err := repo.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return repo.FindByUID(user.UID)
}

// UpdateLastLoginTime updates an last time login
func (repo *UsersRepository) UpdateLastLoginInfo(user *models.User, data *models.User) error {
	updateData := map[string]interface{}{"LastLoginIp": data.LastLoginIp, "LastLoginAt": data.LastLoginAt, "BlockedUntil": nil}
	if err := repo.DB.Model(&user).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UsersRepository) UpdateStatusInfo(user *models.User) error {
	updateData := map[string]interface{}{"Status": user.Status, "BlockedUntil": nil}
	if err := repo.DB.Model(&user).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePasswordAndChallengeName updates password and challenge name
func (repo *UsersRepository) UpdatePasswordAndChallengeName(user *models.User, data *models.User) error {
	updateData := map[string]interface{}{"Password": data.Password, "ChallengeName": data.ChallengeName}
	if err := repo.DB.Model(&user).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// CountByUserGroupID count users by user group id
func (repo *UsersRepository) CountByUserGroupID(userGroupID uint64) (uint64, error) {
	var user models.User
	var count uint64
	if err := repo.DB.Where("user_group_id = ?", userGroupID).Model(&user).Count(&count).
		Error; err != nil {
		return count, err
	}
	return count, nil
}

// GetList returns list by passed params
func (repo *UsersRepository) GetList(params *list_params.ListParams) ([]*models.User, error) {
	var users []*models.User
	adapter := adapters.NewGorm(repo.DB)
	err := adapter.LoadList(&users, params, "users")
	return users, err
}

func (copy UsersRepository) WrapContext(db *gorm.DB) *UsersRepository {
	copy.DB = db
	return &copy
}

func (repo *UsersRepository) Delete(user *models.User) error {
	if err := repo.DB.Delete(user).Error; err != nil {
		return err
	}
	return nil
}

// Process a search query and return SQL conditions(only for columns) and values for them
func (repo *UsersRepository) processSearchQueryParam(searchQuery string, columns []string) (conditions []string, values []interface{}) {
	rawValues := strings.Split(searchQuery, " ")
	conditions = make([]string, 0, len(columns)*len(rawValues))
	values = make([]interface{}, 0, len(columns)*len(rawValues))

	for _, rawValue := range rawValues {
		value := "%" + rawValue + "%"
		for _, column := range columns {
			conditions = append(conditions, fmt.Sprintf("`%s` LIKE ?", column))
			values = append(values, value)
		}
	}

	return conditions, values
}
