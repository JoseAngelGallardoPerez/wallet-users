package models

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Form struct {
	ID uint64 `gorm:"column:id"`

	// Type of a form
	Type string `gorm:"column:type"`

	// List of roles as json
	InitiatorRoleNames string `gorm:"column:initiator_role_names"`

	// Owner of an entity which we want to update
	OwnerRoleNames string `gorm:"column:owner_role_names"`

	// Form as json configuration. See forms.FormConfig
	Form string `gorm:"column:form"`
}

// Decodes json to a list of roles
func (s *Form) InitiatorRoleNamesAsList() ([]string, error) {
	roles := make([]string, 0)
	if s.InitiatorRoleNames == "" {
		return roles, nil
	}

	if err := json.Unmarshal([]byte(s.InitiatorRoleNames), &roles); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal InitiatorRoleNames: %s", s.InitiatorRoleNames)
	}

	return roles, nil
}

// Decodes json to a list of roles
func (s *Form) OwnerRoleNamesAsList() ([]string, error) {
	roles := make([]string, 0)
	if s.OwnerRoleNames == "" {
		return roles, nil
	}

	if err := json.Unmarshal([]byte(s.OwnerRoleNames), &roles); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal OwnerRoleNames: %s", s.OwnerRoleNames)
	}

	return roles, nil
}
