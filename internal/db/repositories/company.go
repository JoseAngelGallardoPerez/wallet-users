package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{
		db,
	}
}

func (this *CompanyRepository) Create(company *models.Company) (*models.Company, error) {
	if err := this.db.Create(company).Error; err != nil {
		return nil, err
	}

	return this.GetByID(company.ID)
}

func (this *CompanyRepository) Update(company *models.Company) (*models.Company, error) {
	if err := this.db.Save(company).Error; err != nil {
		return nil, err
	}

	return this.GetByID(company.ID)
}

func (this *CompanyRepository) GetByID(id uint64) (*models.Company, error) {
	var company models.Company
	if err := this.db.Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}

	return &company, nil
}

func (this *CompanyRepository) GetByName(name string) (*models.Company, error) {
	var company models.Company
	if err := this.db.Where("LOWER(company_name) = ?", name).First(&company).Error; err != nil {
		return nil, err
	}

	return &company, nil
}

func (this *CompanyRepository) SaveAndFindByNames(names []string) ([]*models.Company, error) {
	var companies []*models.Company

	companyNameMap := make(map[string]struct{})
	savedCompanies := make(map[string]struct{})

	if len(names) != len(companies) {
		for _, company := range companies {
			companyNameMap[company.CompanyName] = struct{}{}
		}

		for _, name := range names {
			if _, ok := companyNameMap[name]; !ok {
				if _, ok := savedCompanies[name]; !ok {
					newCompany := models.Company{
						CompanyName: name,
					}

					createdCompany, err := this.Create(&newCompany)
					if err != nil {
						return nil, err
					}

					savedCompanies[name] = struct{}{}
					delete(companyNameMap, name)

					companies = append(companies, createdCompany)
				}
			}
		}
	}

	return companies, nil
}

func (this *CompanyRepository) FindByIDs(ids []uint64) ([]*models.Company, error) {
	var companies []*models.Company
	if err := this.db.Where("id in (?)", ids).Find(&companies).Error; err != nil {
		return nil, err
	}

	return companies, nil
}

func (this *CompanyRepository) FindByID(id uint64) (*models.Company, error) {
	var company models.Company
	if err := this.db.Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}

	return &company, nil
}

func (this *CompanyRepository) DeleteById(id uint64) error {
	if err := this.db.Where("id = ?", id).Delete(&models.Company{}).Error; err != nil {
		return err
	}
	return nil
}

func (this CompanyRepository) WrapContext(db *gorm.DB) *CompanyRepository {
	this.db = db
	return &this
}
