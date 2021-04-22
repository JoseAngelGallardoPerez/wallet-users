package form_conditions

import (
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-users/internal/db/models"
)

const (
	defaultUserClass = "useDefaultClass"
)

type Condition interface {
	// Apply an additional condition to the user model
	Apply(user *models.User, params interface{}) error

	// Returns unique key to identify a condition
	Key() string
}

type ConditionRegistry struct {
	conditions map[string]Condition
}

func NewConditionRegistry() *ConditionRegistry {
	return &ConditionRegistry{conditions: make(map[string]Condition)}
}

// Register a new condition in the factory
func (s *ConditionRegistry) Register(condition Condition) error {
	_, ok := s.conditions[condition.Key()]
	if ok {
		return errors.Errorf("cannot register condition. Key `%s` already exists", condition.Key())
	}

	s.conditions[condition.Key()] = condition
	return nil
}

// Return a registered condition by key
func (s *ConditionRegistry) Make(key string) (Condition, error) {
	condition, ok := s.conditions[key]
	if !ok {
		return nil, errors.Errorf("cannot find condition `%s`", key)
	}

	return condition, nil
}
