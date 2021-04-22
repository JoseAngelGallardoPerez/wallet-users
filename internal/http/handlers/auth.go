package handlers

import (
	"net/http"
	"strings"
	"time"

	errors "github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/forms"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/http/serializers"
	httpAuth "github.com/Confialink/wallet-users/internal/http/services/auth"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/accounts"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/Confialink/wallet-users/internal/services/users"
	"github.com/Confialink/wallet-users/internal/validators"
)

type AuthService struct {
	Repository              repositories.RepositoryInterface
	ResponseService         responses.ResponseHandler
	notificationsService    *notifications.Notifications
	TmpTokensService        *auth.TemporaryTokens
	AuthBlocker             *auth.Blocker
	UserService             *users.UserService
	Logger                  log15.Logger
	PasswordService         *services.Password
	tokenService            *auth.TokenService
	sysSettings             *syssettings.SysSettings
	authService             *auth.Auth
	confirmationCodeService *users.ConfirmationCode
	userForm                *forms.User
	userLoaderService       *users.UserLoaderService
	signUpResponse          *httpAuth.SignUpResponse
	accountsService         *accounts.AccountsService
}

func NewAuthService(
	repository repositories.RepositoryInterface,
	responseService responses.ResponseHandler,
	notificationsService *notifications.Notifications,
	tmpTokensService *auth.TemporaryTokens,
	authBlocker *auth.Blocker,
	userService *users.UserService,
	logger log15.Logger,
	passwordService *services.Password,
	tokenService *auth.TokenService,
	sysSettings *syssettings.SysSettings,
	authService *auth.Auth,
	confirmationCodeService *users.ConfirmationCode,
	userForm *forms.User,
	userLoaderService *users.UserLoaderService,
	signUpResponse *httpAuth.SignUpResponse,
	accountsService *accounts.AccountsService,
) *AuthService {
	return &AuthService{
		Repository:              repository,
		ResponseService:         responseService,
		notificationsService:    notificationsService,
		TmpTokensService:        tmpTokensService,
		AuthBlocker:             authBlocker,
		UserService:             userService,
		Logger:                  logger,
		PasswordService:         passwordService,
		tokenService:            tokenService,
		sysSettings:             sysSettings,
		authService:             authService,
		confirmationCodeService: confirmationCodeService,
		userForm:                userForm,
		userLoaderService:       userLoaderService,
		signUpResponse:          signUpResponse,
		accountsService:         accountsService,
	}
}

// AuthWithBeforeSignIn called before default SignIn
type AuthWithBeforeSignIn interface {
	BeforeSignIn(*gin.Context, *models.User) *responses.Error
}

// AuthWithAfterSignIn called after default SignIn
type AuthWithAfterSignIn interface {
	AfterSignIn(*gin.Context, *models.User) *responses.Error
}

// MeHandler returns current user
func (srv *AuthService) MeHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "MeHandler")

	accessToken, ok := ctx.Get("AccessToken")
	if !ok {
		// Returns a "401 StatusUnauthorized" response
		srv.ResponseService.Error(ctx, responses.Unauthorized, "Access token not found.")
		return
	}

	user, err := srv.Repository.GetUsersRepository().FindUserByTokenAndSubject(accessToken.(string), auth.ClaimAccessSub)
	if err != nil {
		logger.Error("failed to find user by access token", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CannotFindUserByAccessToken, "Can't find user.")
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

// MeLimitedHandler returns current user by X-Tmp-Auth token
func (srv *AuthService) MeLimitedHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "MeLimitedHandler")

	currentUser := GetCurrentUser(ctx)
	if currentUser == nil {
		srv.ResponseService.NotFound(ctx)
		return
	}

	user, err := srv.Repository.GetUsersRepository().FindByUID(currentUser.UID)
	if err != nil {
		logger.Error("failed to find user by uid", "error", err)
		srv.ResponseService.NotFound(ctx)
		return
	}

	serialized := serializers.NewGetUser(user).Serialize()

	// Returns a "200 OK" response
	srv.ResponseService.OkResponse(ctx, serialized)
	return
}

// BeforeSignIn called before SignInHandler
func (srv *AuthService) BeforeSignIn(ctx *gin.Context, user *models.User) *responses.Error {
	var ip = ctx.ClientIP()

	srv.AuthBlocker.LoadSettings()

	if err := srv.AuthBlocker.CheckIP(ip); err != nil {
		e := responses.NewCommonError().ApplyCode(responses.CodeIpIsBlocked)
		return e
	}

	if err := srv.AuthBlocker.CheckUser(user.Email); err != nil {
		return err
	}

	return nil
}

// SignInHandler login a user
func (srv *AuthService) SignInHandler(ctx *gin.Context) {
	var ip = ctx.ClientIP()

	// Checks if the query entry is valid
	validator := validators.LoginValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	if prehook, ok := interface{}(srv).(AuthWithBeforeSignIn); ok {
		if e := prehook.BeforeSignIn(ctx, &validator.UserModel); e != nil {
			r := responses.NewResponse().SetStatus(http.StatusForbidden).AddError(e)
			ctx.AbortWithStatusJSON(r.Status, r)
			return
		}
	}

	user, err := srv.Repository.GetUsersRepository().FindByEmailOrPhoneNumber(validator.UserModel.Email)
	if err != nil {
		srv.AuthBlocker.AddIPFailAttempt(ip)
		// Returns a "401 StatusUnauthorized" response
		srv.ResponseService.Error(ctx, responses.CodeInvalidUsernamePassword, "Invalid username or password.")
		return
	}

	res, errResp := srv.authService.LoginUser(user, validator.UserModel, nil, ip)
	if errResp != nil {
		srv.ResponseService.SetError(ctx, errResp)
		return
	}

	if posthook, ok := interface{}(srv).(AuthWithAfterSignIn); ok {
		if e := posthook.AfterSignIn(ctx, user); e != nil {
			r := responses.NewResponse().SetStatus(http.StatusUnauthorized).AddError(e)
			ctx.AbortWithStatusJSON(r.Status, r)
			return
		}
	}

	srv.ResponseService.SuccessResponse(ctx, http.StatusOK, res)
}

// AfterSignIn called after SignInHandler
func (srv *AuthService) AfterSignIn(ctx *gin.Context, user *models.User) *responses.Error {
	logger := srv.Logger.New("action", "AfterSignIn")

	var ip = ctx.ClientIP()
	var err error

	// set last time login
	lastLoginAt := time.Now()
	info := &models.User{LastLoginIp: ip, LastLoginAt: &lastLoginAt}
	err = srv.Repository.GetUsersRepository().UpdateLastLoginInfo(user, info)
	if err != nil {
		logger.Error("failed to update last login info", "error", err)
	}

	// insert record into accesslog
	err = srv.Repository.GetAccesLogRepository().Create(&models.AccessLog{UID: user.UID, IP: ip})
	if err != nil {
		logger.Error("failed to insert record into accesslog", "error", err)
	}

	srv.AuthBlocker.ClearAllOldFailAttempts(user)

	return nil
}

// RefreshHandler will take in a valid refresh token and return new tokens.
func (srv *AuthService) RefreshHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "RefreshHandler")

	accessToken := ctx.Request.Header.Get("Authorization")
	if len(accessToken) < 8 || !strings.EqualFold(accessToken[0:7], "Bearer ") {
		// Returns a "401 StatusUnauthorized" response
		srv.ResponseService.Error(ctx, responses.Unauthorized, "Access token not found.")
		return
	}

	accessToken = accessToken[7:]

	refreshToken := ctx.Request.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		// Returns a "401 StatusUnauthorized" response
		srv.ResponseService.Error(ctx, responses.Unauthorized, "Refresh token not found.")
		return
	}

	output, err := srv.tokenService.RefreshTokens(accessToken, refreshToken, nil)
	if err != nil {
		logger.Error("failed to refresh token", "error", err)
		// Returns a "401 StatusUnauthorized" response
		srv.ResponseService.Error(ctx, responses.Unauthorized, "")
		return
	}

	// Returns a "200 StatusOK" response
	srv.ResponseService.SuccessResponse(ctx, http.StatusOK, output)
}

// SimpleSignUpHandler registers new a user with less data
func (srv *AuthService) SimpleSignUpHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "SimpleSignUpHandler")
	rawData, err := ctx.GetRawData()
	if err != nil {
		logger.Error("cannot read body", "err", err)
		srv.ResponseService.Error(ctx, responses.CannotCreateUserWithRegistrationRequest, "Can't create user.")
		return
	}

	user, err := srv.userForm.SignUp(rawData)
	if err != nil {
		logger.Error("cannot make a user", "err", err)
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	user, err = srv.UserService.CreateNew(user)
	if err != nil {
		logger.Error("failed to create user by request", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotCreateUserWithRegistrationRequest, "Can't create user.")
		return
	}

	resp, err := srv.signUpResponse.Make(user)
	if err != nil {
		logger.Error("cannot make a correct response", "error", err)
		// We already created a user, so we do not return the error
		// Returns a "201 Created" response
		srv.ResponseService.SuccessResponse(ctx, http.StatusCreated, user)
	}

	if err := srv.accountsService.GenerateAccount(user); err != nil {
		logger.Error("cannot generate account", "error", err, "uid", user.UID)
		// do not return the error in response
	}

	for name, value := range resp.Headers {
		ctx.Header(name, value)
	}
	srv.ResponseService.SuccessResponse(ctx, resp.Status, resp.Data)
}

// SignOutHandler remove the access token
func (srv *AuthService) SignOutHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "SignOutHandler")

	accessToken, ok := ctx.Get("AccessToken")
	if !ok {
		// Returns a "401 StatusUnauthorized" response
		srv.ResponseService.Error(ctx, responses.Unauthorized, "Access token not found.")
		return
	}

	err := srv.tokenService.RevokeToken(accessToken.(string))
	if err != nil {
		logger.Error("failed to sign out", "error", err)
		// Returns a "401 StatusUnauthorized" response
		srv.ResponseService.Error(ctx, responses.Unauthorized, "Can't sign out.")
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

// SignOutDeviceHandler remove the access token and device
func (srv *AuthService) SignOutDeviceHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "SignOutDeviceHandler")

	accessToken, ok := ctx.Get("AccessToken")
	if !ok {
		srv.ResponseService.Error(ctx, responses.Unauthorized, "Access token not found.")
		return
	}

	validator := validators.RemoveDeviceValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	if err := srv.tokenService.RevokeToken(accessToken.(string)); err != nil {
		logger.Error("failed to sign out", "error", err)
		srv.ResponseService.Error(ctx, responses.Unauthorized, "Can't sign out.")
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// ForgotPassword sends to the end user with a confirmation code
func (srv *AuthService) ForgotPassword(ctx *gin.Context) {
	logger := srv.Logger.New("action", "ForgotPassword")

	form := struct {
		Email string `json:"email" binding:"required"`
	}{}

	// Checks if the query entry is valid
	if err := ctx.ShouldBindJSON(&form); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	err := srv.UserService.ForgotPassword(form.Email)
	if err != nil {
		if typedErr, ok := err.(errors.TypedError); ok {
			errors.AddErrors(ctx, typedErr)
			return
		}
		logger.Error("failed to reset password", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CannotForgotPassword, "Can't reset password.")
		return
	}

	// Returns a "200 StatusOK" response
	srv.ResponseService.SuccessResponse(ctx, http.StatusNoContent, nil)
}

// ResetPassword resets password by passing a new password and confirmation code
func (srv *AuthService) ResetPassword(ctx *gin.Context) {
	logger := srv.Logger.New("action", "ResetPassword")

	form := &validators.ResetPassword{}

	// Checks if the query entry is valid
	if err := ctx.ShouldBindJSON(&form); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	err := srv.UserService.ResetPassword(form.NewPassword, form.ConfirmationCode)
	if err != nil {
		logger.Error("failed to reset password", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.InvalidConfirmationCode, "Can't reset password.")
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

// ChangePasswordHandler changes the password for a specified user.
func (srv *AuthService) ChangePasswordHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "ChangePasswordHandler")
	user := ctx.MustGet("_current_user").(*models.User)

	// Checks if the query entry is valid
	validator := validators.ChangePasswordValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	err := srv.PasswordService.UserCheckPassword(
		validator.Model.PreviousPassword,
		user.Password,
	)
	if err != nil {
		logger.Error("failed to change password", "error", err)
		srv.ResponseService.Error(ctx, responses.CodeInvalidPassword, "Invalid password.")
		return
	}

	hash, err := srv.PasswordService.UserHashPassword(validator.Model.ProposedPassword)
	if err != nil {
		logger.Error("failed to create password hash", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CannotCreateHashPassword, "Can't create hash password.")
		return
	}

	user.Password = hash

	if user.ChallengeName != nil && *user.ChallengeName == models.ChallengeNameNewPasswordRequired {
		user.ChallengeName = nil
	}

	err = srv.Repository.GetUsersRepository().UpdatePasswordAndChallengeName(user, user)
	if nil != err {
		logger.Error("failed to update user", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update user.")
		return
	}

	// send notification
	if _, err = srv.notificationsService.PasswordChanged(user.UID); err != nil {
		logger.Error("failed to send notifications for change password user", "error", err)
	}

	// Returns a "204 StatusNoContent" response
	ctx.Status(http.StatusNoContent)
}

// SetPasswordHandler allows user to set password using confirmation code.
func (srv *AuthService) SetPasswordHandler(ctx *gin.Context) {
	logger := srv.Logger.New("action", "SetPasswordHandler")
	user := ctx.MustGet("_current_user").(*models.User)

	// Checks if the query entry is valid
	setPasswordForm := &validators.SetPassword{}
	if err := ctx.ShouldBind(setPasswordForm); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	hash, err := srv.PasswordService.UserHashPassword(setPasswordForm.ProposedPassword)
	if err != nil {
		logger.Error("failed to create password hash", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CannotCreateHashPassword, "Can't create hash password.")
		return
	}

	user.Password = hash

	if user.ChallengeName != nil && *user.ChallengeName == models.ChallengeNameNewPasswordRequired {
		user.ChallengeName = nil
	}

	err = srv.Repository.GetUsersRepository().UpdatePasswordAndChallengeName(user, user)
	if nil != err {
		logger.Error("failed to update user", "error", err)
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update user.")
		return
	}

	// send notification
	if _, err = srv.notificationsService.PasswordChanged(user.UID); err != nil {
		logger.Error("failed to send notifications for change password user", "error", err)
	}

	tmpAuthToken, err := srv.TmpTokensService.Issue(user)
	if err != nil {
		logger.Error("failed to create temporary auth token", "error", err)
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CannotCreateTmpAuthToken, "Can't create temporary auth token.")
		return
	}
	ctx.Header(httpAuth.TmpAuthHeader, tmpAuthToken)
	ctx.Header("Access-Control-Expose-Headers", httpAuth.TmpAuthHeader)

	ctx.Status(http.StatusCreated)
}

func (srv *AuthService) GetConfirmationCodeHandler(ctx *gin.Context) {
	confirmationCode := ctx.MustGet("_confirmation_code").(*models.ConfirmationCode)
	srv.ResponseService.OkResponse(ctx, confirmationCode)
}

func (srv *AuthService) IssueTokensForUserByUID(ctx *gin.Context) {
	logger := srv.Logger.New("action", "IssueTokensForUserByUID")

	currentUser := GetCurrentUser(ctx)
	//TODO: do not compare role name directly
	if currentUser.RoleName != "root" {
		srv.ResponseService.Error(ctx, responses.Forbidden, "Action is not allowed")
		ctx.Abort()
		return
	}

	uid, ok := ctx.Params.Get("uid")
	if !ok {
		logger.Error("uid parameter is not found")
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()
		return
	}

	user, err := srv.Repository.GetUsersRepository().FindByUID(uid)
	if err != nil {
		logger.Error("can't find user", err)
		srv.ResponseService.Error(ctx, responses.NotFound, "")
		return
	}

	res, err := srv.tokenService.IssueTokens(user, nil)
	if err != nil {
		logger.Error("can't issue tokens", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	srv.ResponseService.SuccessResponse(ctx, http.StatusOK, res)
}
