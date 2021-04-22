package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/Confialink/wallet-pkg-utils/csv"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/forms"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/http/serializers"
	"github.com/Confialink/wallet-users/internal/services"
	csvServices "github.com/Confialink/wallet-users/internal/services/csv"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/permissions"
	systemLogs "github.com/Confialink/wallet-users/internal/services/system-logs"
	"github.com/Confialink/wallet-users/internal/services/users"
	"github.com/Confialink/wallet-users/internal/validators"
)

// UsersService contains handler for users API
type UsersService struct {
	Repository              repositories.RepositoryInterface
	ResponseService         responses.ResponseHandler
	SystemLogsService       *systemLogs.SystemLogsService
	notificationsService    *notifications.Notifications
	PermissionsService      *permissions.Permissions
	userProfilesCsvService  *csvServices.Users
	adminProfilesCsvService *csvServices.AdminProfiles
	userCreator             *users.UserService
	params                  *HandlerParams
	logger                  log15.Logger
	PasswordService         *services.Password
	confirmationCodeService *users.ConfirmationCode
	userLoaderService       *users.UserLoaderService
	userForm                *forms.User
}

// NewUsersService return new UsersService
func NewUsersService(
	repository repositories.RepositoryInterface,
	responseService responses.ResponseHandler,
	systemLogsService *systemLogs.SystemLogsService,
	notificationsService *notifications.Notifications,
	permissionsService *permissions.Permissions,
	userProfilesCsvService *csvServices.Users,
	adminProfilesCsvService *csvServices.AdminProfiles,
	userCreator *users.UserService,
	params *HandlerParams,
	logger log15.Logger,
	PasswordService *services.Password,
	confirmationCodeService *users.ConfirmationCode,
	userLoaderService *users.UserLoaderService,
	userForm *forms.User,
) *UsersService {
	return &UsersService{
		repository,
		responseService,
		systemLogsService,
		notificationsService,
		permissionsService,
		userProfilesCsvService,
		adminProfilesCsvService,
		userCreator,
		params,
		logger,
		PasswordService,
		confirmationCodeService,
		userLoaderService,
		userForm,
	}
}

// GetAdminProfilesCsvHandler is handler to download csv file with admin profiles
func (srv *UsersService) GetAdminProfilesCsvHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "GetAdminProfilesCsvHandler")

	var errs []*responses.Error

	params := srv.params.adminProfilesCsv(ctx.Request.URL.RawQuery)
	if ok, errorsList := params.Validate(); !ok {
		for _, err := range errorsList {
			e := responses.NewCommonError().
				ApplyCode(responses.CannotGetAdminProfilesAsCsv).
				SetDetails(err.Error())
			errs = append(errs, e)
		}
		srv.ResponseService.Errors(ctx, http.StatusBadRequest, errs)
		return
	}
	params.AddFilter("roleName", []string{models.RoleAdmin})

	file, err := srv.adminProfilesCsvService.GetFile(params)
	if err != nil {
		logger.Error("Can not get csv file", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotGetAdminProfilesAsCsv, "Can't get csv file")
		return
	}

	if err := csv.Send(file, ctx.Writer); err != nil {
		logger.Error("Can not send csv file", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotSendUserProfilesAsCsv, "Can't send csv file")
		return
	}
}

// GetUserProfilesCsvHandler is handler to download csv file with user profiles
func (srv *UsersService) GetUserProfilesCsvHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "GetUsersCsvHandler")

	var errs []*responses.Error

	params := srv.params.userProfilesCsv(ctx.Request.URL.RawQuery)
	params.AddFilter("role_name", []string{
		models.RoleClient,
	}, list_params.OperatorIn)
	if ok, errorsList := params.Validate(); !ok {
		for _, err := range errorsList {
			e := responses.NewCommonError().
				ApplyCode(responses.CannotGetUserProfilesAsCsv).
				SetDetails(err.Error())
			errs = append(errs, e)
		}
		srv.ResponseService.Errors(ctx, http.StatusBadRequest, errs)
		return
	}

	file, err := srv.userProfilesCsvService.GetFile(params)
	if err != nil {
		logger.Error("Can not get csv file", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotGetUserProfilesAsCsv, "Can't get csv file")
		return
	}

	if err = csv.Send(file, ctx.Writer); err != nil {
		logger.Error("Can not send csv file", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotSendUserProfilesAsCsv, "Can't send csv file")
		return
	}
}

// ListHandler returns the list of users
func (srv *UsersService) ListHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "ListHandler")

	currentUser := GetCurrentUser(ctx)

	limitQuery := ctx.DefaultQuery("limit", "10")
	pageQuery := ctx.DefaultQuery("page", "1")
	params := ctx.Request.URL.Query()

	var adminsRoleIncluded = false

	roles := params["filter[role_name]"]
	if len(roles) > 0 {
		for key, role := range roles {
			// remove root from role names if user is not root
			// only root can see root users
			if role == models.RoleRoot && currentUser.RoleName != models.RoleRoot {
				copy(roles[key:], roles[key+1:])
				roles[len(roles)-1] = ""
				roles = roles[:len(roles)-1]
			}
			if role == models.RoleRoot || role == models.RoleAdmin {
				adminsRoleIncluded = true
			}
		}
	} else {
		adminsRoleIncluded = true
	}

	var hasPerm bool
	if adminsRoleIncluded {
		hasPerm = srv.PermissionsService.CanViewAdminProfile(currentUser.UID)
	} else {
		hasPerm = srv.PermissionsService.CanViewUserProfile(currentUser.UID)
	}

	if !hasPerm {
		srv.ResponseService.Forbidden(ctx)
		return
	}

	query := srv.Repository.GetUsersRepository().Filter(params)

	pagination, err := srv.Repository.GetUsersRepository().Paginate(query, pageQuery, limitQuery, serializers.NewUsers())
	if err != nil {
		logger.Error("сan't load list of user", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of users")
		return
	}

	// Returns a "200 OK" response
	srv.ResponseService.OkResponse(ctx, pagination)
}

// ShortListHandler returns list of users with limited list of fields
func (srv *UsersService) ShortListHandler(ctx *gin.Context) {
	limitQuery := ctx.DefaultQuery("limit", "10")
	pageQuery := ctx.DefaultQuery("page", "1")
	query := srv.Repository.GetUsersRepository().Filter(ctx.Request.URL.Query())

	pagination, err := srv.Repository.GetUsersRepository().Paginate(query, pageQuery, limitQuery, serializers.NewShortUsers())
	if err != nil {
		logger := srv.logger.New("action", "ShortListHandler")
		logger.Error("сan't load list of user", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of users")
		return
	}

	srv.ResponseService.OkResponse(ctx, pagination)
}

// ListContacts returns list of active users with limited list of fields
func (srv *UsersService) ListContacts(ctx *gin.Context) {
	logger := srv.logger.New("action", "ListContacts")

	form := forms.ListContacts{}
	if err := form.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	usersList, err := srv.Repository.GetUsersRepository().FindActiveByPhoneNumbers(form.PhoneNumbers)
	if err != nil {
		logger.Error("сan't load list of user", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of users")
		return
	}

	shortener := serializers.NewShortUsers()
	srv.ResponseService.OkResponse(ctx, shortener.Short(usersList, serializers.ShortContactsFields))
	return
}

// GetHandler returns user by uid
func (srv *UsersService) GetHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "GetHandler")

	user := GetRequestedUser(ctx)
	if user == nil {
		logger.Error("сan't load a user")
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	serializedUser, err := srv.userLoaderService.LoadUserCompletelyAndSerialize(user)
	if err != nil {
		logger.Error("failed to assemble a user", "error", err)
		srv.ResponseService.Error(ctx, responses.InternalError, "")
		return
	}

	// Returns a "200 OK" response
	srv.ResponseService.OkResponse(ctx, serializedUser)
	return
}

// GetShortHandler returns user by uid with limited list of fields
func (srv *UsersService) GetShortHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "GetShortHandler")

	user := GetRequestedUser(ctx)
	if user == nil {
		logger.Error("сan't load a user")
		srv.ResponseService.NotFound(ctx)
		return
	}

	serialized := serializers.NewShortUser(user).Serialize()
	srv.ResponseService.OkResponse(ctx, serialized)
	return
}

// CreateHandler creates new user
func (srv *UsersService) CreateHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "CreateHandler")

	// Checks if the query entry is valid
	validator := validators.CreateUserValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	// Check permissions
	currentUser := GetCurrentUser(ctx)
	if hasPerm := srv.PermissionsService.CanCreateProfile(currentUser.UID, &validator.UserModel); !hasPerm {
		srv.ResponseService.Forbidden(ctx)
		return
	}

	tmpPassword := validator.UserModel.Password

	// Create new user
	createdUser, err := srv.userCreator.Create(&validator.UserModel, true, false, nil)
	if err != nil {
		logger.Error("сan't create a user", "error", err)
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CanNotCreateUser, "Can't create a user")
		return
	}

	if nil != currentUser {
		srv.SystemLogsService.LogCreateUserProfileAsync(createdUser, currentUser.UID)
	}
	// TODO: refactor - use events, move above functionality to the event subscriber
	confirmationCode, err := srv.confirmationCodeService.GenerateSetPasswordCode(createdUser)
	if err != nil {
		logger.Error("unable to generate set_password confirmation code")
		return
	}

	if _, err = srv.notificationsService.ProfileCreated(createdUser.UID, tmpPassword, confirmationCode.Code); nil != err {
		logger.Error("сan't send notification", "error", err)
		return
	}

	// Returns a "201 Created" response
	srv.ResponseService.SuccessResponse(ctx, http.StatusCreated, validator.UserModel)
}

// UpdateHandler updates a user
func (srv *UsersService) UpdateHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "UpdateHandler")
	user := GetRequestedUser(ctx)
	if user == nil {
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	rawData, err := ctx.GetRawData()
	if err != nil {
		logger.Error("cannot read body", "err", err)
		srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update user.")
		return
	}

	currentUser := GetCurrentUser(ctx)
	if currentUser.UID == user.UID ||
		currentUser.RoleName == "root" ||
		currentUser.RoleName == "admin" {

		err = srv.userForm.Update(user, currentUser, rawData)
		if err != nil {
			logger.Error("cannot update a user", "err", err)
			srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
			return
		}

		old, err := srv.Repository.GetUsersRepository().FindByUID(user.UID)
		if err != nil {
			logger.Error("cannot found user", "err", err)
			srv.ResponseService.NotFound(ctx)
			return
		}

		err = srv.userLoaderService.LoadUserCompletely(old)
		if err != nil {
			logger.Error("cannot load user", "err", err)
			srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update a user")
			return
		}

		tx := srv.Repository.GetUsersRepository().DB.Begin()
		err = srv.userCreator.Update(user, tx)
		if err != nil {
			tx.Rollback()
			// Returns a "400 StatusBadRequest" response
			srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update a user")
			return
		}

		if currentUser.UID != user.UID &&
			(currentUser.RoleName == "admin" || currentUser.RoleName == "root") {
			srv.SystemLogsService.LogModifyUserProfileAsync(old, user, currentUser.UID)
		}

		tx.Commit()
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

// PatchHandler updates a user fields
func (srv *UsersService) PatchHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "PatchHandler")
	user := GetRequestedUser(ctx)
	if user == nil {
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	// Checks if the query entry is valid
	form := &validators.PatchUser{}
	if err := ctx.ShouldBindJSON(form); err != nil {
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	currentUser := GetCurrentUser(ctx)
	if currentUser.UID == user.UID ||
		currentUser.RoleName == "root" ||
		currentUser.RoleName == "admin" {

		if form.FirstName != nil {
			user.FirstName = *form.FirstName
		}
		if form.LastName != nil {
			user.LastName = *form.LastName
		}
		if form.Nickname != nil {
			user.Nickname = *form.Nickname
		}

		repo := srv.Repository.GetUsersRepository()
		old, err := repo.FindByUID(user.UID)
		if err != nil {
			logger.Error("cannot find user", "err", err)
			srv.ResponseService.NotFound(ctx)
			return
		}

		_, err = repo.Update(user)
		if err != nil {
			logger.Error("cannot update user", "err", err)
			srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update a user")
			return
		}

		if currentUser.UID != user.UID &&
			(currentUser.RoleName == "admin" || currentUser.RoleName == "root") {
			srv.SystemLogsService.LogModifyUserProfileAsync(old, user, currentUser.UID)
		}
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

func (srv *UsersService) ResetPasswordHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "UpdateHandler")
	user := GetRequestedUser(ctx)
	if user == nil {
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	currentUser := GetCurrentUser(ctx)

	if currentUser.UID == user.UID ||
		currentUser.RoleName == "root" ||
		currentUser.RoleName == "admin" {

		formResetPassword := forms.ResetPassword{}
		if err := formResetPassword.BindJSON(ctx); err != nil {
			// Returns a "422 StatusUnprocessableEntity" response
			srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
			return
		}

		err := srv.userCreator.SetPassword(user, formResetPassword.NewPassword)
		if err != nil {
			// Returns a "400 StatusBadRequest" response
			srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update a user")
			return
		}

		if _, err := srv.notificationsService.PasswordChanged(user.UID); err != nil {
			logger.Error("Can't send notification.", ctx, err)
			return
		}
	}
}

// UnblockHandler unblock users
func (srv *UsersService) UnblockHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "UpdateHandler")
	// Checks if the query entry is valid
	validator := validators.UnblockUsersValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	// Unblock users
	for _, item := range validator.Data {
		err := srv.unblockUser(item.UID)
		if nil != err {
			logger.Error("failed unblock user", "error", err)
			srv.ResponseService.Error(ctx, responses.CanNotUnblockUser, "Can't unblock user")
			return
		}
	}
	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

// ExportHandler exports users to csv
func (srv *UsersService) ExportHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "ExportHandler")
	csv := services.Csv{}

	var users []*models.User

	query := srv.Repository.GetUsersRepository().Filter(ctx.Request.URL.Query())
	if err := query.Find(&users).Error; err != nil {
		logger.Error("failed export users", "error", err)
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CanNotExportUsers, "Can't export users")
		return
	}

	b := csv.UsersToCsv(users)

	fileName := fmt.Sprintf("users-%s.csv", time.RFC3339)
	ctx.Writer.Header().Set("Content-Type", "text/csv")
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)
	ctx.Writer.Write(b.Bytes())
}

// ImportHandler import users from csv
func (srv *UsersService) ImportHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "ImportHandler")
	csv := services.Csv{}

	file, _, err := ctx.Request.FormFile("file")
	defer file.Close()

	if nil != err {
		logger.Error("failed import users", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CanNotImportUsers, "Can't import users")
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		logger.Error("failed import users", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CanNotImportUsers, "Can't import users")
		return
	}

	users, err := csv.CsvToUsers(buf)
	if nil != err {
		logger.Error("failed import users", "error", err)
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CanNotImportUsers, "Can't import users")
		return
	}

	var errors []*responses.Error
	for _, user := range users {
		// Create new user
		_, err = srv.userCreator.Create(user, true, true, nil)
		if err != nil {
			logger.Error("failed create user via import csv", "error", err)
			e := responses.NewCommonError().
				ApplyCode(responses.CannotUnblockIp).
				SetDetails(fmt.Sprintf("Can't create user with uid `%s`", user.UID))
			errors = append(errors, e)
		}
	}

	if len(errors) > 0 {
		srv.ResponseService.Errors(ctx, http.StatusInternalServerError, errors)
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

func (srv *UsersService) GenerateNewPhoneCode(ctx *gin.Context) {
	logger := srv.logger.New("action", "GenerateNewPhoneCode")

	user := GetRequestedUser(ctx)
	if user == nil {
		srv.ResponseService.NotFound(ctx)
		return
	}
	// generate new confirmation code
	code, err := srv.confirmationCodeService.CreateNewVerificationCode(user, models.ConfirmationCodeSubjectPhoneVerificationCode)
	if err != nil {
		logger.Error("failed create new phone verification code", "error", err)
		srv.ResponseService.Error(ctx, responses.CanNotGeneratePhoneVerificationCode, "Can't generate phone verification code")
		return
	}
	// send notification
	if _, err = srv.notificationsService.VerifyPhone(user.UID, code.Code); err != nil {
		logger.Error("сan't send notification for phone verification", "error", err)
		srv.ResponseService.Error(ctx, responses.CanNotGeneratePhoneVerificationCode, "Can't send notification for phone verification")
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

func (srv *UsersService) CheckPhoneCode(ctx *gin.Context) {
	logger := srv.logger.New("action", "CheckPhoneCode")

	user := GetRequestedUser(ctx)
	if user == nil {
		srv.ResponseService.NotFound(ctx)
		return
	}

	phoneCode := &struct {
		Code string `json:"code"`
	}{}

	if err := ctx.Bind(phoneCode); err != nil {
		logger.Error("сan't bind json", "error", err)
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	user.IsPhoneConfirmed = true
	_, err := srv.userCreator.CheckPhoneCode(phoneCode.Code, user)
	if err != nil {
		logger.Error("phone verification failed", "error", err)
		srv.ResponseService.Error(ctx, responses.InvalidConfirmationCode, "")
		return
	}
	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

func (srv *UsersService) GenerateNewEmailCode(ctx *gin.Context) {
	logger := srv.logger.New("action", "GenerateNewEmailCode")

	user := GetRequestedUser(ctx)
	if user == nil {
		srv.ResponseService.NotFound(ctx)
		return
	}

	code, err := srv.confirmationCodeService.CreateNewVerificationCode(user, models.ConfirmationCodeSubjectEmailVerificationCode)
	if err != nil {
		logger.Error("failed create new email verification code", "error", err)
		srv.ResponseService.Error(ctx, responses.CanNotGenerateEmailVerificationCode, "Can't generate email verification code")
		return
	}

	// send notification
	if _, err = srv.notificationsService.VerifyEmail(user.UID, code.Code); nil != err {
		logger.Error("сan't send notification for email verification", "error", err)
		srv.ResponseService.Error(ctx, responses.CanNotGenerateEmailVerificationCode, "Can't send notification for email verification")
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

func (srv *UsersService) CheckEmailCode(ctx *gin.Context) {
	logger := srv.logger.New("action", "CheckEmailCode")

	user := GetRequestedUser(ctx)
	if user == nil {
		srv.ResponseService.NotFound(ctx)
		return
	}

	emailCode := &struct {
		Code string `json:"code"`
	}{}

	if err := ctx.Bind(emailCode); err != nil {
		logger.Error("сan't bind json", "error", err)
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	_, err := srv.userCreator.CheckEmailCode(emailCode.Code, user)
	if err != nil {
		logger.Error("email verification failed", "error", err)
		srv.ResponseService.Error(ctx, responses.InvalidConfirmationCode, "")
		return
	}
	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

//
// Helper functions
//

// unblockUser is helper function for unblock user
func (srv *UsersService) unblockUser(uid string) error {

	user, err := srv.Repository.GetUsersRepository().FindByUID(uid)
	if err != nil {
		return err
	}

	user.ClearBlockedUntil()

	_, err = srv.Repository.GetUsersRepository().Save(user)
	if nil != err {
		return err
	}

	return nil
}
