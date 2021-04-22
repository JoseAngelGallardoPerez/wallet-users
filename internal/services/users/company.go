package users

import (
	"errors"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
)

type CompanyService struct {
	companyRepository *repositories.CompanyRepository
}

func NewCompanyService(
	companyRepository *repositories.CompanyRepository,
) *CompanyService {
	return &CompanyService{companyRepository}
}

// UpdateCompanyDetails creates/updates a new db record or remove an existing record if need
func (s *CompanyService) UpdateCompanyDetails(user *models.User, tx *gorm.DB) error {
	repo := s.companyRepository.WrapContext(tx)
	details := user.CompanyDetails

	if details.ID != 0 && user.CompanyID != nil {
		if details.ID != *user.CompanyID {
			return errors.New("this company does not belong to the user")
		}
		if _, err := repo.Update(&details); err != nil {
			return err
		}
	} else if details.CompanyName != "" ||
		details.CompanyType != "" || details.CompanyRole != "" ||
		details.DirectorFirstName != "" || details.DirectorLastName != "" {

		company, err := repo.Create(&details)
		if err != nil {
			return err
		}
		user.CompanyID = &company.ID
	} else if user.CompanyID != nil {
		if err := repo.DeleteById(*user.CompanyID); err != nil {
			return err
		}

		user.CompanyID = nil
	}

	return nil
}
