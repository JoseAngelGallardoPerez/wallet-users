package serializers

import (
	"github.com/Confialink/wallet-pkg-model_serializer"

	"github.com/Confialink/wallet-users/internal/db/models"
)

var ShortUserFields = []interface{}{"UID", "Email", "Username", "FirstName", "LastName", "Nickname", "Status", "LastLoginAt", "ParentId"}
var ShortContactsFields = []interface{}{"UID", "PhoneNumber"}

type shortUser struct {
	user *models.User
}

func NewShortUser(user *models.User) *shortUser {
	return &shortUser{user}
}

func (s *shortUser) Serialize() map[string]interface{} {
	return model_serializer.Serialize(s.user, ShortUserFields)
}
