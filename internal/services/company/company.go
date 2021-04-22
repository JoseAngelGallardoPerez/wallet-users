package company

import (
	"strings"

	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
)

type CompanyService struct {
	companyRepository repositories.CompanyRepository
	logger            log15.Logger
}

func NewCompanyService(
	companyRepository repositories.CompanyRepository,
	logger log15.Logger,
) *CompanyService {
	return &CompanyService{
		companyRepository,
		logger,
	}
}

func (this *CompanyService) Create(company *models.Company) (*models.Company, error) {
	return this.companyRepository.Create(company)
}

func (this *CompanyService) Update(company *models.Company) (*models.Company, error) {
	return this.companyRepository.Update(company)
}

func (this *CompanyService) GetByID(id uint64) (*models.Company, error) {
	return this.companyRepository.GetByID(id)
}

func (this *CompanyService) GetByName(name string) (*models.Company, error) {
	return this.companyRepository.GetByName(strings.ToLower(name))
}
