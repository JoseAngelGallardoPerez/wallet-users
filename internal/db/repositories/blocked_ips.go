package repositories

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/jinzhu/gorm"
)

// BlockedIpsRepositoryHandler is interface for repository functionality that ought to be implemented manually.
type BlockedIpsRepositoryHandler interface {
	FindByUID(uid string) (*models.BlockedIp, error)
	Create(blockedIp *models.BlockedIp) error
}

// BlockedIpsRepository is userGroup group repository for CRUD operations.
type BlockedIpsRepository struct {
	DB *gorm.DB
}

func NewBlockedIpsRepository(DB *gorm.DB) *BlockedIpsRepository {
	return &BlockedIpsRepository{DB: DB}
}

// Paginate returns a new Pagination instance.
func (repo *BlockedIpsRepository) Paginate(query *gorm.DB, pageQuery string, limitQuery string) (*Pagination, error) {
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

	var ips []*models.BlockedIp
	var count int

	go countItems(query, ips, done, &count)

	if err := query.Raw("SELECT `id`, INET_NTOA(`ip`) as ip, `created_at`, `blocked_until` FROM blocked_ips").
		Limit(p.Limit).Offset(p.Offset).Scan(&ips).Error; err != nil {
		return nil, err
	}
	<-done

	p.TotalRecord = count
	p.Items = ips
	p.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	return p, nil
}

// Filter apply request params to the builder instance.
func (repo *BlockedIpsRepository) Filter(params url.Values) *gorm.DB {
	db := repo.DB
	db = repo.applySort(db, params)
	return db
}

func (repo *BlockedIpsRepository) applySort(query *gorm.DB, params url.Values) *gorm.DB {
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

// Create creates new record
func (repo *BlockedIpsRepository) Create(ip string, blockedUntil time.Time) error {
	if err := repo.DB.Exec("INSERT INTO blocked_ips (`ip`, `created_at`, `blocked_until`) VALUES (INET_ATON(?), CURRENT_TIMESTAMP(), ?)", ip, blockedUntil).
		Error; err != nil {
		return err
	}
	return nil
}

// Delete an existing ip
func (repo *BlockedIpsRepository) Delete(blockedIp *models.BlockedIp) error {
	if err := repo.DB.Delete(blockedIp).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BlockedIpsRepository) FindByIp(ip string) (*models.BlockedIp, error) {
	blockedIp := &models.BlockedIp{}
	if err := repo.DB.Where("ip = INET_ATON(?)", ip).First(&blockedIp).Error; err != nil {
		return nil, err
	}
	return blockedIp, nil
}
