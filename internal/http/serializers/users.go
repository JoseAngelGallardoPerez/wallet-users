package serializers

import (
	"github.com/Confialink/wallet-pkg-model_serializer"
	"github.com/Confialink/wallet-users/internal/db/models"
)

type shortUsers struct {
}

func NewShortUsers() *shortUsers {
	return &shortUsers{}
}

func (s *shortUsers) Serialize(users []*models.User) []interface{} {
	return s.Short(users, ShortUserFields)
}

func (s *shortUsers) Short(users []*models.User, fields []interface{}) []interface{} {
	res := make([]interface{}, len(users), len(users))
	for key, user := range users {
		res[key] = model_serializer.Serialize(user, fields)
	}
	return res
}

type users struct {
}

func NewUsers() *users {
	return &users{}
}

func (s *users) Serialize(users []*models.User) []interface{} {
	res := make([]interface{}, len(users), len(users))
	for key, user := range users {
		res[key] = *user
	}
	return res
}
