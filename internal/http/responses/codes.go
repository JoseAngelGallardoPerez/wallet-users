package responses

import "net/http"

const (
	Forbidden                               = "FORBIDDEN"
	NotFound                                = "NOT_FOUND"
	Unauthorized                            = "UNAUTHORIZED"
	InternalError                           = "INTERNAL_ERROR"
	ActionCannotBePerformed                 = "ACTION_CANNOT_BE_PERFORMED"
	CodeIpIsBlocked                         = "USERS_IP_IS_BLOCKED"
	CodeUserIsBlocked                       = "USERS_USER_IS_BLOCKED"
	CodeUserIsPending                       = "USERS_USER_IS_PENDING"
	CodeAuthUserIsNotFound                  = "USERS_AUTH_USER_IS_NOT_FOUND"
	CodeRegistrationIsNotConfirmed          = "USERS_REGISTRATION_IS_NOT_CONFIRMED"
	CodeUserIsLocked                        = "USERS_USER_IS_LOCKED"
	CodeUserIsNotActive                     = "USERS_USER_IS_NOT_ACTIVE"
	CodeUserIsDormant                       = "USERS_USER_IS_DORMANT"
	CodeResetPasswordIsNotAllowed           = "RESET_PASSWORD_IS_NOT_ALLOWED"
	CodeConfirmRegistrationLink             = "USERS_CONFIRM_REGISTRATION_LINK"
	CodeInvalidUsernamePassword             = "USERS_INVALID_USERNAME_OR_PASSWORD"
	CodeInvalidPin                          = "INVALID_PIN"
	CodeInvalidPassword                     = "USERS_INVALID_PASSWORD"
	BadCollectionParams                     = "BAD_COLLECTION_PARAMS"
	CannotRetrieveCollection                = "CANNOT_RETRIEVE_COLLECTION"
	CannotRetrieveCollectionCount           = "CANNOT_RETRIEVE_COLLECTION_COUNT"
	CannotCreatePagination                  = "CANNOT_CREATE_PAGINATION"
	CanNotApproveUser                       = "CANNOT_APPROVE_USER"
	CanNotUpdateUser                        = "CANNOT_UPDATE_USER"
	CanNotCreateUser                        = "CANNOT_CREATE_USER"
	CodeInvalidSecurityQuestionAnswer       = "USERS_INVALID_SECURITY_QUESTION_ANSWER"
	CannotCreateDevice                      = "CANNOT_CREATE_DEVICE"
	DeviceNotFound                          = "DEVICE_NOT_FOUND"
	UserGroupAlreadyExists                  = "USER_GROUP_ALREADY_EXISTS"
	CanNotCreateUserGroup                   = "CANNOT_CREATE_USER_GROUP"
	CanNotUpdateUserGroup                   = "CANNOT_UPDATE_USER_GROUP"
	CanNotDeleteUserGroup                   = "CANNOT_DELETE_USER_GROUP"
	InvalidIdForUserGroup                   = "INVALID_ID_FOR_USER_GROUP"
	CannotUpdateSecurityQuestionAnswer      = "CANNOT_UPDATE_SECURITY_QUESTION_ANSWER"
	CanNotUpdateUserSettings                = "CANNOT_UPDATE_USER_SETTINGS"
	CanNotUnblockUser                       = "CANNOT_UNBLOCK_USER"
	CanNotExportUsers                       = "CANNOT_EXPORT_USERS"
	CanNotImportUsers                       = "CANNOT_IMPORT_USERS"
	CannotGetUserProfilesAsCsv              = "CANNOT_GET_USER_PROFILES_AS_CSV"
	CannotGetAdminProfilesAsCsv             = "CANNOT_GET_ADMIN_PROFILES_AS_CSV"
	CannotSendUserProfilesAsCsv             = "CANNOT_SEND_USER_PROFILES_AS_CSV"
	CannotCreateUserWithRegistrationRequest = "CANNOT_CREATE_USER_WITH_REGISTRATION_REQUEST"
	CannotCreateTmpAuthToken                = "CANNOT_CREATE_TMP_AUTH_TOKEN"
	CannotFindUserByEmail                   = "CANNOT_FIND_USER_BY_EMAIL"
	CannotFindUserByAccessToken             = "CANNOT_FIND_USER_BY_ACCESS_TOKEN"
	CannotForgotPassword                    = "CANNOT_FORGOT_PASSWORD"
	CannotFindSecurityQuestionsAnswerByUser = "CANNOT_FIND_SECURITY_QUESTIONS_ANSWERS_BY_USER"
	CannotConfirmForgotPassword             = "CANNOT_CONFIRM_FORGOT_PASSWORD"
	InvalidConfirmationCode                 = "INVALID_CONFIRMATION_CODE"
	ConfirmationCodeIsInvalid               = "CONFIRMATION_CODE_IS_INVALID"
	CannotCreateHashPassword                = "CANNOT_CREATE_HASH_PASSWORD"
	CannotChangePassword                    = "CANNOT_CHANGE_PASSWORD"
	CanNotAddVerification                   = "CANNOT_ADD_VERIFICATION"
	CanNotAddInvite                         = "CANNOT_ADD_INVITE"
	CanNotUpdateVerification                = "CANNOT_UPDATE_VERIFICATION"
	CanNotCreateVerificationRequest         = "CANNOT_CREATE_VERIFICATION_REQUEST"
	CanNotApproveVerificationRequest        = "CANNOT_APPROVE_VERIFICATION_REQUEST"
	CanNotCancelVerificationRequest         = "CANNOT_CANCEL_VERIFICATION_REQUEST"
	VerificationNotFound                    = "VERIFICATION_NOT_FOUND"
	MaxVerificationFiles                    = "MAX_VERIFICATION_FILES"
	CanNotGeneratePhoneVerificationCode     = "CANNOT_GENERATE_PHONE_VERIFICATION_CODE"
	CanNotGenerateEmailVerificationCode     = "CANNOT_GENERATE_EMAIL_VERIFICATION_CODE"
	PhoneNumberIsNotConfirmed               = "PHONE_NUMBER_IS_NOT_CONFIRMED"

	UnprocessableEntity       = "UNPROCESSABLE_ENTITY"
	DocumentTypeOneOf         = "DOCUMENT_TYPE_ONE_OF"
	CannotUnblockIp           = "CANNOT_UNBLOCK_IP"
	StatusOneOf               = "STATUS_ONE_OF"
	EmailAlreadyExists        = "EMAIL_ALREADY_EXISTS"
	UnsupportedRole           = "UNSUPPORTED_ROLE"
	PhoneAlreadyExists        = "PHONE_ALREADY_EXISTS"
	UsernameAlreadyExists     = "USERNAME_ALREADY_EXISTS"
	PhoneNumber               = "PHONE_NUMBER"
	GDRP                      = "GDRP"
	SpecialCharacterRequired  = "SPECIAL_CHARACTER_REQUIRED"
	NumberRequired            = "NUMBER_REQUIRED"
	UppercaseLetterRequired   = "UPPERCASE_LETTER_REQUIRED"
	LowercaseLetterRequired   = "LOWERCASE_LETTER_REQUIRED"
	UnknownEmailOrPhoneNumber = "UNKNOWN_EMAIL_OR_PHONE_NUMBER"

	MaintenanceMode = "MAINTENANCE_MODE"
)

var statusCodes = map[string]int{
	Forbidden:                               http.StatusForbidden,
	NotFound:                                http.StatusNotFound,
	Unauthorized:                            http.StatusUnauthorized,
	InternalError:                           http.StatusInternalServerError,
	ActionCannotBePerformed:                 http.StatusUnauthorized,
	CodeIpIsBlocked:                         http.StatusUnauthorized,
	CodeUserIsBlocked:                       http.StatusForbidden,
	CodeUserIsPending:                       http.StatusUnauthorized,
	CodeAuthUserIsNotFound:                  http.StatusUnauthorized,
	CodeRegistrationIsNotConfirmed:          http.StatusUnauthorized,
	CodeUserIsLocked:                        http.StatusLocked,
	CodeUserIsNotActive:                     http.StatusForbidden,
	CodeResetPasswordIsNotAllowed:           http.StatusForbidden,
	CodeUserIsDormant:                       http.StatusUnauthorized,
	CodeConfirmRegistrationLink:             http.StatusUnauthorized,
	CodeInvalidUsernamePassword:             http.StatusUnauthorized,
	CodeInvalidPin:                          http.StatusUnauthorized,
	CodeInvalidPassword:                     http.StatusUnprocessableEntity,
	BadCollectionParams:                     http.StatusBadRequest,
	CannotRetrieveCollection:                http.StatusBadRequest,
	CannotRetrieveCollectionCount:           http.StatusBadRequest,
	CannotCreatePagination:                  http.StatusBadRequest,
	CanNotApproveUser:                       http.StatusBadRequest,
	CanNotUpdateUser:                        http.StatusBadRequest,
	CanNotCreateUser:                        http.StatusInternalServerError,
	CodeInvalidSecurityQuestionAnswer:       http.StatusUnauthorized,
	CannotCreateDevice:                      http.StatusInternalServerError,
	DeviceNotFound:                          http.StatusBadRequest,
	UserGroupAlreadyExists:                  http.StatusConflict,
	CanNotCreateUserGroup:                   http.StatusInternalServerError,
	CanNotUpdateUserGroup:                   http.StatusInternalServerError,
	CanNotDeleteUserGroup:                   http.StatusInternalServerError,
	InvalidIdForUserGroup:                   http.StatusBadRequest,
	CannotUpdateSecurityQuestionAnswer:      http.StatusInternalServerError,
	CanNotUpdateUserSettings:                http.StatusInternalServerError,
	CanNotUnblockUser:                       http.StatusBadRequest,
	CanNotExportUsers:                       http.StatusInternalServerError,
	CanNotImportUsers:                       http.StatusInternalServerError,
	CannotGetUserProfilesAsCsv:              http.StatusBadRequest,
	CannotGetAdminProfilesAsCsv:             http.StatusBadRequest,
	CannotSendUserProfilesAsCsv:             http.StatusInternalServerError,
	CannotCreateUserWithRegistrationRequest: http.StatusInternalServerError,
	CannotCreateTmpAuthToken:                http.StatusInternalServerError,
	CannotFindUserByEmail:                   http.StatusBadRequest,
	CannotFindUserByAccessToken:             http.StatusBadRequest,
	CannotForgotPassword:                    http.StatusBadRequest,
	CannotConfirmForgotPassword:             http.StatusBadRequest,
	CannotFindSecurityQuestionsAnswerByUser: http.StatusBadRequest,
	CannotCreateHashPassword:                http.StatusInternalServerError,
	CannotChangePassword:                    http.StatusInternalServerError,
	InvalidConfirmationCode:                 http.StatusBadRequest,
	ConfirmationCodeIsInvalid:               http.StatusBadRequest,
	CanNotAddVerification:                   http.StatusInternalServerError,
	CanNotAddInvite:                         http.StatusInternalServerError,
	CanNotUpdateVerification:                http.StatusInternalServerError,
	CanNotCreateVerificationRequest:         http.StatusBadRequest,
	VerificationNotFound:                    http.StatusNotFound,
	MaxVerificationFiles:                    http.StatusBadRequest,
	PhoneNumberIsNotConfirmed:               http.StatusForbidden,

	UnprocessableEntity:      http.StatusUnprocessableEntity,
	DocumentTypeOneOf:        http.StatusUnprocessableEntity,
	StatusOneOf:              http.StatusUnprocessableEntity,
	EmailAlreadyExists:       http.StatusUnprocessableEntity,
	PhoneNumber:              http.StatusUnprocessableEntity,
	GDRP:                     http.StatusUnprocessableEntity,
	CannotUnblockIp:          http.StatusInternalServerError,
	SpecialCharacterRequired: http.StatusUnprocessableEntity,
	NumberRequired:           http.StatusUnprocessableEntity,
	UppercaseLetterRequired:  http.StatusUnprocessableEntity,
	LowercaseLetterRequired:  http.StatusUnprocessableEntity,

	MaintenanceMode: http.StatusForbidden,
}
