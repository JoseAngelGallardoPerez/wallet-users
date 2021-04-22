package repositories

import (
	"github.com/jinzhu/gorm"
)

// RepositoryInterface is interface for repository functionality
type RepositoryInterface interface {
	GetUsersRepository() *UsersRepository
	GetUserGroupsRepository() *UserGroupsRepository
	GetSecurityQuestionRepository() *SecurityQuestionRepository
	GetSecurityQuestionsAnswerRepository() *SecurityQuestionsAnswerRepository
	GetAccesLogRepository() *AccesLogRepository
	GetBlockedIpsRepository() *BlockedIpsRepository
	GetFailAuthAttemptRepository() *FailAuthAttemptRepository
	GetVerificationRepository() *VerificationRepository
	GetInvitesRepository() *InvitesRepository
	GetCompanyRepository() *CompanyRepository
}

// Repository is user repository for CRUD operations.
type Repository struct {
	DB *gorm.DB
}

// Pagination is the abstract pagination
type Pagination struct {
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Items       interface{} `json:"items"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
}

// NewRepository creates new user repository
func NewRepository(db *gorm.DB) RepositoryInterface {
	return &Repository{db}
}

// GetUsersRepository gets the repository for a users
func (repo *Repository) GetUsersRepository() *UsersRepository {
	return &UsersRepository{repo.DB}
}

// GetUserGroupsRepository gets the repository for a user groups
func (repo *Repository) GetUserGroupsRepository() *UserGroupsRepository {
	return &UserGroupsRepository{repo.DB}
}

// GetSecurityQuestionRepository gets the repository for an security questions
func (repo *Repository) GetSecurityQuestionRepository() *SecurityQuestionRepository {
	return &SecurityQuestionRepository{repo.DB}
}

// GetSecurityQuestionsAnswerRepository gets the repository for an security questions
func (repo *Repository) GetSecurityQuestionsAnswerRepository() *SecurityQuestionsAnswerRepository {
	return &SecurityQuestionsAnswerRepository{repo.DB}
}

// GetAccesLogRepository gets the repository for an settings
func (repo *Repository) GetAccesLogRepository() *AccesLogRepository {
	return &AccesLogRepository{repo.DB}
}

// GetBlockedIpsRepository gets the repository for an blocked ips
func (repo *Repository) GetBlockedIpsRepository() *BlockedIpsRepository {
	return &BlockedIpsRepository{repo.DB}
}

func (repo *Repository) GetFailAuthAttemptRepository() *FailAuthAttemptRepository {
	return &FailAuthAttemptRepository{repo.DB}
}

func (repo *Repository) GetVerificationRepository() *VerificationRepository {
	return &VerificationRepository{repo.DB}
}

func (repo *Repository) GetInvitesRepository() *InvitesRepository {
	return &InvitesRepository{repo.DB}
}

// GetCompanyRepository gets the repository for an company
func (repo *Repository) GetCompanyRepository() *CompanyRepository {
	return &CompanyRepository{repo.DB}
}

// countItems gets how many records for a query
func countItems(query *gorm.DB, users interface{}, done chan bool, count *int) {
	query.Model(users).Count(count)
	done <- true
}
