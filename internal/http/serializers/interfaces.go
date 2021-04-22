package serializers

import (
	"github.com/Confialink/wallet-users/internal/db/models"
)

type UsersCollectionSerializer interface {
	Serialize([]*models.User) []interface{}
}
