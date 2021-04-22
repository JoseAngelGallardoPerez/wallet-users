package form_conditions

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

type DefaultUserClass struct {
	sysSettings *syssettings.SysSettings
}

func NewDefaultUserClass(sysSettings *syssettings.SysSettings) *DefaultUserClass {
	return &DefaultUserClass{sysSettings: sysSettings}
}

func (s *DefaultUserClass) Apply(user *models.User, params interface{}) error {
	classId, err := s.sysSettings.GetDefaultUserClassByRole(user.RoleName)
	if err != nil {
		return errors.Wrap(err, "cannot receive default user class")
	}

	user.ClassId = json.Number(*classId)

	return nil
}

func (s *DefaultUserClass) Key() string {
	return defaultUserClass
}
