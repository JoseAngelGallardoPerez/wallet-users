package users

import (
	"strconv"

	"github.com/Confialink/wallet-users/internal/db/models"
)

func ToTypedValue(value, valueType string) interface{} {

	switch valueType {
	case models.AttributeTypeString:
		return value
	case models.AttributeTypeBool:
		return value == "true"
	case models.AttributeTypeInt:
		res, _ := strconv.ParseInt(value, 10, 64)
		return res
	case models.AttributeTypeFloat:
		res, _ := strconv.ParseFloat(value, 64)
		return res
	}

	return value // default string
}
