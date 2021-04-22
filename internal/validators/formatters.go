package validators

import (
	"fmt"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/go-playground/validator/v10"

	"github.com/Confialink/wallet-users/internal/http/responses"
)

func formatters() map[string]*errors.ValidationErrorFormatter {
	return map[string]*errors.ValidationErrorFormatter{
		"documenttypeoneof": {
			Code: responses.DocumentTypeOneOf,
			TitleFunc: func(_ validator.FieldError, formattedField string) string {
				return "Unknown type of document"
			},
		},
		"statusoneof": {
			Code: responses.StatusOneOf,
			TitleFunc: func(_ validator.FieldError, formattedField string) string {
				return "Unknown status"
			},
		},
		"uniqueEmail": {
			Code: responses.EmailAlreadyExists,
			TitleFunc: func(field validator.FieldError, formattedField string) string {
				return fmt.Sprintf("`%s` already exists", field.Value())
			},
		},
		"uniqueUsername": {
			Code: responses.UsernameAlreadyExists,
			TitleFunc: func(field validator.FieldError, formattedField string) string {
				return fmt.Sprintf("`%s` already exists", field.Value())
			},
		},
		"phonenumber": {
			Code: responses.PhoneNumber,
			TitleFunc: func(_ validator.FieldError, formattedField string) string {
				return "Invalid format"
			},
		},
		"specialCharacterRequired": {
			Code: responses.SpecialCharacterRequired,
			TitleFunc: func(field validator.FieldError, formattedField string) string {
				return fmt.Sprintf(`%s must contain at least one special character`, field.Field())
			},
		},
		"numberRequired": {
			Code: responses.NumberRequired,
			TitleFunc: func(f validator.FieldError, formattedField string) string {
				return fmt.Sprintf(`%s must contain at least one number`, f.Field())
			},
		},
		"uppercaseLetterRequired": {
			Code: responses.UppercaseLetterRequired,
			TitleFunc: func(f validator.FieldError, formattedField string) string {
				return fmt.Sprintf(`%s must contain at least one uppercase letter`, f.Field())
			},
		},
		"lowercaseLetterRequired": {
			Code: responses.LowercaseLetterRequired,
			TitleFunc: func(f validator.FieldError, formattedField string) string {
				return fmt.Sprintf(`%s must contain at least one lowercase letter`, f.Field())
			},
		},
		"uniquePhoneNumber": {
			Code: responses.PhoneAlreadyExists,
			TitleFunc: func(f validator.FieldError, formattedField string) string {
				return "Phone number already exists on the platform"
			},
		},
	}
}
