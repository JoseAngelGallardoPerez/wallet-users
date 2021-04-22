package validators

import (
	"encoding/json"
	"time"

	"github.com/Confialink/wallet-pkg-utils/pointer"
	"github.com/gin-gonic/gin/binding"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type RootValidator struct {
	Email    string `binding:"required,email"`
	Password string `binding:"required,min=8,max=255,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`

	UserModel models.User `json:"-"`
}

func (s *RootValidator) Call() error {
	if err := binding.Validator.ValidateStruct(s); err != nil {
		return err
	}

	s.UserModel.Email = s.Email
	s.UserModel.Password = s.Password

	s.UserModel.IsCorporate = pointer.ToBool(false)
	s.UserModel.RoleName = models.RoleRoot
	s.UserModel.ClassId = json.Number("1")
	s.UserModel.Status = models.StatusActive
	s.UserModel.LastActedAt = time.Now()

	return nil
}
