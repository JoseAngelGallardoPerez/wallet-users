package validators

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"reflect"
	"regexp"
	"time"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/go-playground/validator/v10"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

var (
	phoneNumberPattern = regexp.MustCompile(`^\+[\d]{7,15}$`) // +XXX XXX XXX XXX
	// One of !"#$%&'\()*+,-./:;<=>?@[\]^_{|}~
	// Reference: https://owasp.org/www-community/password-special-characters
	atLeastOneSpecialCharacterPattern = regexp.MustCompile(`[!"#$%&'\\()*+,-./:;<=>?@[\]^_{|}~]`)
	atLeastOneNumberPattern           = regexp.MustCompile(`[0-9]`)
	atLeastOneUppercaseLetterPattern  = regexp.MustCompile(`[A-Z]`)
	atLeastOneLowercaseLetterPattern  = regexp.MustCompile(`[a-z]`)
	usernamePattern                   = regexp.MustCompile(`^[a-zA-Z0-9_\.'\-\s]+$`)
)

// Interface is a validator interface
type Interface interface {
	// Validate a structure
	Struct(current interface{}) error
}

func Register(
	usersRepo *repositories.UsersRepository,
	securityQuestionsRepo *repositories.SecurityQuestionRepository,
	userGroupsRepo *repositories.UserGroupsRepository,
	sysSettings *syssettings.SysSettings,
	passwordService *services.Password,
	v *validator.Validate,
) {
	if err := v.RegisterValidation("roleoneof", roleOneOf); err != nil {
		panic("cannot add validation roleoneof")
	}

	if err := v.RegisterValidation("documenttypeoneof", documentTypeOneOf); err != nil {
		panic("cannot add validation documenttypeoneof")
	}
	if err := v.RegisterValidation("statusoneof", statusOneOf); err != nil {
		panic("cannot add validation statusoneof")
	}
	if err := v.RegisterValidation("uniqueEmail", uniqueEmail(usersRepo)); err != nil {
		panic("cannot add validation uniqueEmail")
	}
	if err := v.RegisterValidation("uniqueUsername", noUsernameExists(usersRepo)); err != nil {
		panic("cannot add validation uniqueUsername")
	}
	if err := v.RegisterValidation("uniquePhoneNumber", uniquePhoneNumber(usersRepo)); err != nil {
		panic("cannot add validation uniquePhoneNumber")
	}
	if err := v.RegisterValidation("usernameChars", usernameChars); err != nil {
		panic("cannot add validation usernameChars")
	}
	if err := v.RegisterValidation("securityQuestionExists", securityQuestionExists(securityQuestionsRepo)); err != nil {
		panic("cannot add validation securityQuestionExists")
	}
	if err := v.RegisterValidation("phonenumber", isValidPhoneNumber); err != nil {
		panic("cannot add validation phonenumber")
	}
	if err := v.RegisterValidation("gdpr", isGdprAccepted); err != nil {
		panic("cannot add validation gdpr")
	}
	if err := v.RegisterValidation("specialCharacterRequired", specialCharacterRequired); err != nil {
		panic("cannot add validation specialCharacterRequired")
	}
	if err := v.RegisterValidation("numberRequired", numberRequired); err != nil {
		panic("cannot add validation numberRequired")
	}
	if err := v.RegisterValidation("uppercaseLetterRequired", uppercaseLetterRequired); err != nil {
		panic("cannot add validation uppercaseLetterRequired")
	}
	if err := v.RegisterValidation("lowercaseLetterRequired", lowercaseLetterRequired); err != nil {
		panic("cannot add validation lowercaseLetterRequired")
	}
	if err := v.RegisterValidation("noGroupNameExists", noGroupNameExists(userGroupsRepo)); err != nil {
		panic("cannot add validation noGroupNameExists")
	}
	if err := v.RegisterValidation("verificationTypeOneOf", verificationTypeOneOf); err != nil {
		panic("cannot add validation verificationTypeOneOf")
	}

	if err := v.RegisterValidation("existCountry", existCountry); err != nil {
		panic("cannot add validation existCountry")
	}
	if err := v.RegisterValidation("dayBeforeNow", dayBeforeNow); err != nil {
		panic("cannot add validation dayBeforeNow")
	}

	v.RegisterStructValidation(updateUserStructLevelValidator(usersRepo, passwordService), UpdateUserValidator{})

	errors.SetFormatters(formatters())
}

// ValidationError is the abstract validation error model
type ValidationError struct {
	Errors map[string]interface{} `json:"errors"`
}

// roleOneOf checks if role is valid
func roleOneOf(fl validator.FieldLevel) bool {
	// TODO: implement if need: receive roles from DB
	return true
}

// roleOneOf checks if document is valid
func documentTypeOneOf(fl validator.FieldLevel) bool {
	if documentType, ok := fl.Field().Interface().(sql.NullString); ok {
		if len(documentType.String) > 0 && !models.IsValidDocumentType(documentType.String) {
			return false
		}
	}
	return true
}

// statusOneOf checks if status is valid
func statusOneOf(fl validator.FieldLevel) bool {
	if status, ok := fl.Field().Interface().(string); ok {
		validStatuses := models.GetAvailableStatuses()
		if len(status) > 0 && !validStatuses[status] {
			return false
		}
	}
	return true
}

func verificationTypeOneOf(fl validator.FieldLevel) bool {
	if verificationType, ok := fl.Field().Interface().(string); ok {
		validVerificationTypes := models.GetAvailableVerificationTypes()
		if len(verificationType) > 0 && !validVerificationTypes[verificationType] {
			return false
		}
	}
	return true
}

// uniqueEmail validates if email does not exist in a given repository.
// It receives a param - name of struct field to get UID user.
// example:
// `binding:"uniqueEmail"`
// `binding:"uniqueEmail=Uid"`
func uniqueEmail(repo *repositories.UsersRepository) validator.Func {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(string); ok {
			// 1. Find the user by email in database
			user, err := repo.FindByEmail(value)
			if err != nil && err != gorm.ErrRecordNotFound {
				return false
			}

			if user == nil {
				return true
			}

			if len(fl.Param()) != 0 {
				field := reflect.Indirect(fl.Parent()).FieldByName(fl.Param())
				if !field.IsValid() {
					return false
				}

				if user.UID == field.String() {
					return true
				}
			}

			return false
		}
		return true
	}
}

// noUsernameExists validates if username does not exist in a given repository
// It receives a param - name of struct field to get UID user.
// example:
// `binding:"uniqueUsername"`
// `binding:"uniqueUsername=Uid"`
func noUsernameExists(repo *repositories.UsersRepository) validator.Func {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(string); ok {
			// 1. Find the user by email in database
			user, err := repo.FindByUsername(value)
			if err != nil && err != gorm.ErrRecordNotFound {
				return false
			}

			if user == nil {
				return true
			}

			if len(fl.Param()) != 0 {
				field := reflect.Indirect(fl.Parent()).FieldByName(fl.Param())
				if !field.IsValid() {
					return false
				}

				if user.UID == field.String() {
					return true
				}
			}

			return false
		}
		return true
	}
}

// uniquePhoneNumber validates if phone_number does not exist in a given repository
// It receives a param - name of struct field to get UID user.
// example:
// `binding:"uniquePhoneNumber"`
// `binding:"uniquePhoneNumber=Uid"`
func uniquePhoneNumber(repo *repositories.UsersRepository) validator.Func {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(string); ok {
			// 1. Find the user by phone_number in database
			user, err := repo.FindByPhoneNumber(value)
			if err != nil && err != gorm.ErrRecordNotFound {
				return false
			}

			if user == nil {
				return true
			}

			if len(fl.Param()) != 0 {
				field := reflect.Indirect(fl.Parent()).FieldByName(fl.Param())
				if !field.IsValid() {
					return false
				}

				if user.UID == field.String() {
					return true
				}
			}

			return false
		}
		return true
	}
}

// usernameChars validates if username includes only allowed chars
func usernameChars(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if len(value) > 0 {
			return usernamePattern.MatchString(value)
		}

	}
	return true
}

// securityQuestionExists validates if a security question exists in a given repository
func securityQuestionExists(repo *repositories.SecurityQuestionRepository) validator.Func {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(uint64); ok {
			questions, err := repo.FindBySqids([]uint64{value})
			if err != nil || len(questions) == 0 {
				return false
			}
		}
		return true
	}
}

// specialCharacterRequired validates if string contains at least one special character
func specialCharacterRequired(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if len(value) > 0 && !atLeastOneSpecialCharacterPattern.MatchString(value) {
			return false
		}
	}
	return true
}

// numberRequired validates if string contains at least one number
func numberRequired(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if len(value) > 0 && !atLeastOneNumberPattern.MatchString(value) {
			return false
		}
	}
	return true
}

// uppercaseLetterRequired validates if string contains at least one uppercase letter
func uppercaseLetterRequired(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if len(value) > 0 && !atLeastOneUppercaseLetterPattern.MatchString(value) {
			return false
		}
	}
	return true
}

// lowercaseLetterRequired validates if string contains at least one lowercase letter
func lowercaseLetterRequired(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if len(value) > 0 && !atLeastOneLowercaseLetterPattern.MatchString(value) {
			return false
		}
	}
	return true
}

func isValidPhoneNumber(fl validator.FieldLevel) bool {
	if phonenumber, ok := fl.Field().Interface().(string); ok {
		if len(phonenumber) > 0 && !phoneNumberPattern.MatchString(phonenumber) {
			return false
		}
	}
	return true
}

func isGdprAccepted(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(bool)
}

func noGroupNameExists(repo *repositories.UserGroupsRepository) validator.Func {
	return func(fl validator.FieldLevel) bool {
		val, valOk := fl.Field().Interface().(string)
		vld, vldOk := fl.Top().Interface().(*UpdateUserGroupValidator)
		if valOk && vldOk && vld.UserGroupModel.Name != val {
			_, err := repo.FindByName(val)
			if nil == err {
				return false
			}
		}
		return true
	}
}

func existCountry(fl validator.FieldLevel) bool {
	top := fl.Field().Interface().(string)
	switch top {
	case
		"AF", "AX", "AL", "DZ", "AS", "AD", "AO", "AI", "AQ", "AG", "AR", "AM", "AW", "AU", "AT", "AZ", "BS", "BH", "BD", "BB", "BY", "BE",
		"BZ", "BJ", "BM", "BT", "BO", "BQ", "BA", "BW", "BV", "BR", "IO", "BN", "BG", "BF", "BI", "CV", "KH", "CM", "CA", "KY", "CF", "TD",
		"CL", "CN", "CX", "CC", "CO", "KM", "CD", "CG", "CK", "CR", "CI", "HR", "CU", "CW", "CY", "CZ", "DK", "DJ", "DM", "DO", "EC", "EG",
		"SV", "GQ", "ER", "EE", "ET", "FK", "FO", "FJ", "FI", "FR", "GF", "PF", "TF", "GA", "GM", "GE", "DE", "GH", "GI", "GR", "GL", "GD",
		"GP", "GU", "GT", "GG", "GN", "GW", "GY", "HT", "HM", "VA", "HN", "HK", "HU", "IS", "IN", "ID", "IR", "IQ", "IE", "IM", "IL", "IT",
		"JM", "JP", "JE", "JO", "KZ", "KE", "KI", "KP", "KR", "KW", "KG", "LA", "LV", "LB", "LS", "LR", "LY", "LI", "LT", "LU", "MO", "MK",
		"MG", "MW", "MY", "MV", "ML", "MT", "MH", "MQ", "MR", "MU", "YT", "MX", "FM", "MD", "MC", "MN", "ME", "MS", "MA", "MZ", "MM", "NA",
		"NR", "NP", "NL", "NC", "NZ", "NI", "NE", "NG", "NU", "NF", "MP", "NO", "OM", "PK", "PW", "PS", "PA", "PG", "PY", "PE", "PH", "PN",
		"PL", "PT", "PR", "QA", "RE", "RO", "RU", "RW", "BL", "SH", "KN", "LC", "MF", "PM", "VC", "WS", "SM", "ST", "SA", "SN", "RS", "SC",
		"SL", "SG", "SX", "SK", "SI", "SB", "SO", "ZA", "GS", "SS", "ES", "LK", "SD", "SR", "SJ", "SZ", "SE", "CH", "SY", "TW", "TJ", "TZ",
		"TH", "TL", "TG", "TK", "TO", "TT", "TN", "TR", "TM", "TC", "TV", "UG", "UA", "AE", "GB", "UM", "US", "UY", "UZ", "VU", "VE", "VN",
		"VG", "VI", "WF", "EH", "YE", "ZM", "ZW":
		return true
	}
	return false
}

func updateUserStructLevelValidator(repo *repositories.UsersRepository, passwordService *services.Password) validator.StructLevelFunc {
	return func(sl validator.StructLevel) {
		s := sl.Current().Interface().(UpdateUserValidator)

		user, _ := repo.FindByUID(s.Data.UID)

		// ensure user not exist with given email if email has been changed
		if user.Email != s.Data.Email {
			// 1. Find the user by email in database
			_, errDb := repo.FindByEmail(s.Data.Email)

			// Report error if user exist in database
			if errDb == nil {
				sl.ReportError(reflect.ValueOf(s.Data.Email), "Email", "Email", "uniqueEmail", "")
			}
		}

		// ensure user not exist with given username if username has been changed
		if user.Username != s.Data.Username {
			// 1. Find the user by username in database
			_, errDb := repo.FindByUsername(s.Data.Username)

			// Report error if user exist in database
			if errDb == nil {
				sl.ReportError(reflect.ValueOf(s.Data.Username), "Username", "Username", "uniqueUsername", "")
			}
		}

		if s.Data.PreviousPassword != nil && s.Data.Password != nil && *s.Data.PreviousPassword != "" && *s.Data.Password != "" {
			if s.CurrentUser == nil || s.CurrentUser.UID == s.Data.UID ||
				(s.CurrentUser != nil && s.CurrentUser.RoleName != models.RoleRoot && s.CurrentUser.RoleName != models.RoleAdmin) {
				err := passwordService.UserCheckPassword(*s.Data.PreviousPassword, user.Password)
				if err != nil {
					sl.ReportError(reflect.ValueOf(*s.Data.PreviousPassword), "PreviousPassword", "PreviousPassword", "invalid", "")
				}
			}

			if *s.Data.Password != *s.Data.ConfirmPassword {
				sl.ReportError(reflect.ValueOf(*s.Data.Password), "ConfirmPassword", "ConfirmPassword", "invalid", "")
			}

			if !atLeastOneSpecialCharacterPattern.MatchString(*s.Data.Password) {
				sl.ReportError(reflect.ValueOf(*s.Data.Password), "Password", "Password", "specialCharacterRequired", "")
			}

			if !atLeastOneNumberPattern.MatchString(*s.Data.Password) {
				sl.ReportError(reflect.ValueOf(*s.Data.Password), "Password", "Password", "numberRequired", "")
			}

			if !atLeastOneUppercaseLetterPattern.MatchString(*s.Data.Password) {
				sl.ReportError(reflect.ValueOf(*s.Data.Password), "Password", "Password", "uppercaseLetterRequired", "")
			}

			if !atLeastOneLowercaseLetterPattern.MatchString(*s.Data.Password) {
				sl.ReportError(reflect.ValueOf(*s.Data.Password), "Password", "Password", "lowercaseLetterRequired", "")
			}
		}
	}
}

// Validates entered date. It must be before current time.
func dayBeforeNow(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		datetime, err := time.Parse("2006-01-02", value)
		if err != nil {
			return false
		}

		return datetime.Before(time.Now())
	}
	return true
}
