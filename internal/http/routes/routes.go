package routes

import (
	"github.com/Confialink/wallet-users/internal/authentication"
	"github.com/Confialink/wallet-users/internal/config"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/handlers"
	"github.com/Confialink/wallet-users/internal/http/middlewares"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/permissions"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/Confialink/wallet-users/internal/services/users"
	"github.com/Confialink/wallet-users/internal/version"

	"github.com/Confialink/wallet-pkg-env_mods"
	errorsPkg "github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"net/http"
)

func Api(
	cfg *config.Configuration,
	logger log15.Logger,

	usersHandler *handlers.UsersService,
	staffsHandler *handlers.StaffsService,
	authHandler *handlers.AuthService,
	userGroupsHandler *handlers.UserGroupsService,
	securityQuestionHandler *handlers.SecurityQuestionService,
	securityQuestionsAnswersHandler *handlers.SecurityQuestionAnswerService,
	blockedIpsHandler *handlers.BlockedIpsService,
	verificationsHandler *handlers.VerificationHandler,
	invitesHandler *handlers.InvitesHandler,

	responseService responses.ResponseHandler,
	usersRepository *repositories.UsersRepository,
	tmpTokens *auth.TemporaryTokens,
	sysSettings *syssettings.SysSettings,
	confirmationCodeService *users.ConfirmationCode,
	permissionsService *permissions.Permissions,
) *gin.Engine {
	// Retrieve config options.
	ginMode := env_mods.GetMode(cfg.GetServer().GetEnv())
	gin.SetMode(ginMode)

	// Creates a gin router with default middleware:
	r := gin.New()

	r.GET("/users/health-check", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/users/build", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.BuildInfo)
	})

	mwUserFromAccessToken := middlewares.UserFromAccessToken(responseService, usersRepository)

	r.Use(
		middlewares.CorsMiddleware(),
		gin.Recovery(),
		gin.Logger(),
	)

	apiGroup := r.Group("users")
	apiGroup.Use(
		errorsPkg.ErrorHandler(logger),
	)

	mwTmpAuth := middlewares.TmpAuthByToken(tmpTokens, usersRepository, logger)
	mwMaintenance := middlewares.Maintenance(sysSettings)

	// TODO: refactor - get rid of hardcoded layers, create a set of composable elements
	mwAdminOrRoot := middlewares.AdminOrRoot()
	mwOwnerOrAdminOrRoot := middlewares.OwnerOrAdminOrRoot()
	mwRequestedUser := middlewares.RequestedUser(usersRepository, responseService)
	mwPermissionsService := middlewares.NewPermissionsMiddleware(permissionsService, responseService)

	/*
	 |---------------------------------------------------
	 | Private router group
	 |---------------------------------------------------
	*/
	privateGroup := apiGroup.Group("/private")
	{
		v1Group := privateGroup.Group("/v1", authentication.Middleware(logger.New("middleware", "Auth")))
		{
			// POST /users/private/v1/list-contacts
			v1Group.POST("/list-contacts", usersHandler.ListContacts)

			usersGroup := v1Group.Group("/users")
			{
				// GET /users/private/v1/users
				usersGroup.GET("/", mwAdminOrRoot, usersHandler.ListHandler)
				// GET /users/private/v1/users/:uid
				usersGroup.GET("/:uid", mwAdminOrRoot, mwRequestedUser, mwPermissionsService.CanViewProfile(), usersHandler.GetHandler)
				// POST /users/private/v1/users
				usersGroup.POST("/", mwAdminOrRoot, usersHandler.CreateHandler)
				// PUT /users/private/v1/users/:uid
				usersGroup.PUT("/:uid", mwOwnerOrAdminOrRoot, mwRequestedUser, mwPermissionsService.CanUpdateProfile(), usersHandler.UpdateHandler)
				usersGroup.PATCH("/:uid", mwOwnerOrAdminOrRoot, mwRequestedUser, mwPermissionsService.CanUpdateProfile(), usersHandler.PatchHandler)
				// PUT /users/private/v1/users/:uid/reset-password
				usersGroup.PUT("/:uid/reset-password", mwOwnerOrAdminOrRoot, mwRequestedUser, mwPermissionsService.CanUpdateProfile(), usersHandler.ResetPasswordHandler)
				// POST /users/private/v1/users/unblock
				usersGroup.POST("/unblock", mwAdminOrRoot, usersHandler.UnblockHandler)
			}

			staffsGroup := v1Group.Group("/staffs")
			{
				staffsGroup.GET("/", staffsHandler.ListHandler)
				staffsGroup.GET("/:uid", mwRequestedUser, mwPermissionsService.CanViewProfile(), staffsHandler.GetHandler)
				staffsGroup.POST("/", staffsHandler.CreateHandler)
				staffsGroup.PUT("/:uid", mwRequestedUser, mwPermissionsService.CanUpdateProfile(), staffsHandler.UpdateHandler)
			}

			shortUsersGroup := v1Group.Group("/short-users", mwAdminOrRoot, mwPermissionsService.CanViewShortUserProfiles())
			{
				// GET /users/private/v1/short-users
				shortUsersGroup.GET("/", usersHandler.ShortListHandler)
				// GET /users/private/v1/short-users/:uid
				shortUsersGroup.GET("/:uid", mwRequestedUser, usersHandler.GetShortHandler)
			}

			exportGroup := v1Group.Group("/export")
			{
				// GET /users/private/v1/export/users
				exportGroup.GET("/users", mwAdminOrRoot, usersHandler.ExportHandler)
				exportGroup.GET("/user-profiles", mwAdminOrRoot, mwPermissionsService.CanViewClientProfile(), usersHandler.GetUserProfilesCsvHandler)
				exportGroup.GET("/admin-profiles", mwAdminOrRoot, mwPermissionsService.CanViewAdminProfile(), usersHandler.GetAdminProfilesCsvHandler)
			}

			authGroup := v1Group.Group("auth")
			{
				// GET /users/private/v1/auth/me
				authGroup.GET("/me", authHandler.MeHandler)
				// DELETE /users/private/v1/auth/logout
				authGroup.DELETE("/logout", authHandler.SignOutHandler)
				// POST /users/private/v1/auth/logout-device
				authGroup.POST("/logout-device", authHandler.SignOutDeviceHandler)
				// POST /users/private/v1/auth/change_password
				authGroup.POST("/change_password", mwUserFromAccessToken, authHandler.ChangePasswordHandler)
				// POST /users/private/v1/auth/root/issue-tokens-for-user-by-uid/:uid
				authGroup.POST("/root/issue-tokens-for-user-by-uid/:uid", authHandler.IssueTokensForUserByUID)

				mwCurrentUserAsRequestedUser := middlewares.CurrentUserAsRequested(usersRepository, responseService)
				// POST /users/private/v1/auth/generate-new-phone-code
				authGroup.POST("/generate-new-phone-code", mwCurrentUserAsRequestedUser, usersHandler.GenerateNewPhoneCode)
				// PUT /users/private/v1/auth/check-phone-code
				authGroup.PUT("/check-phone-code", mwCurrentUserAsRequestedUser, usersHandler.CheckPhoneCode)
				// POST /users/private/v1/auth/generate-new-email-code
				authGroup.POST("/generate-new-email-code", mwCurrentUserAsRequestedUser, usersHandler.GenerateNewEmailCode)
				// PUT /users/private/v1/auth/check-email-code
				authGroup.PUT("/check-email-code", mwCurrentUserAsRequestedUser, usersHandler.CheckEmailCode)
			}

			userGroupsGroup := v1Group.Group("/user-groups", mwAdminOrRoot)
			{
				// GET /users/private/v1/user-groups
				userGroupsGroup.GET("", userGroupsHandler.ListHandler)
				// GET /users/private/v1/user-groups/:uid
				userGroupsGroup.GET("/:id", userGroupsHandler.GetHandler)
				// POST /users/private/v1/user-groups
				userGroupsGroup.POST("", mwPermissionsService.CanCreateSettings(), userGroupsHandler.CreateHandler)
				// PUT /users/private/v1/user-groups/:id
				userGroupsGroup.PUT("/:id", mwPermissionsService.CanModifySettings(), userGroupsHandler.UpdateHandler)
				// DELETE /users/private/v1/user-groups/:id
				userGroupsGroup.DELETE("/:id", mwPermissionsService.CanRemoveSettings(), userGroupsHandler.DeleteHandler)
			}

			securityQuestionsAnswersGroup := v1Group.Group("/security-questions-answers")
			{
				// GET /users/private/v1/security-questions-answers/:uid
				securityQuestionsAnswersGroup.GET("/:uid", mwOwnerOrAdminOrRoot, mwRequestedUser, securityQuestionsAnswersHandler.GetSecurityQuestionAnswersHandler)
				// PUT /users/private/v1/security-questions-answers/:uid
				securityQuestionsAnswersGroup.PUT("/:uid", mwOwnerOrAdminOrRoot, mwRequestedUser, securityQuestionsAnswersHandler.UpdateSecurityQuestionAnswersHandler)
			}

			blockedIpsGroup := v1Group.Group("/blocked-ips")
			{
				// GET /users/private/v1/blocked-ips
				blockedIpsGroup.GET("", mwAdminOrRoot, blockedIpsHandler.ListHandler)
				// POST /users/private/v1/blocked-ips/unblock
				blockedIpsGroup.POST("/unblock", mwAdminOrRoot, blockedIpsHandler.UnblockIpsHandler)
			}

			verificationsGroup := v1Group.Group("/verifications")
			{
				verificationsGroup.GET("/list/:uid", mwOwnerOrAdminOrRoot, verificationsHandler.ListHandler)

				verificationsGroup.POST("", verificationsHandler.CreateHandler)
				verificationsGroup.GET("/verify/:id", verificationsHandler.RequestHandler)
				verificationsGroup.GET("/approve/:id", mwAdminOrRoot, verificationsHandler.ApproveHandler)
				verificationsGroup.GET("/cancel/:id", mwAdminOrRoot, verificationsHandler.CancelHandler)
			}

			invitesGroup := v1Group.Group("/invites")
			{
				invitesGroup.GET("/count", invitesHandler.CountHandler)
				invitesGroup.POST("", invitesHandler.CreateHandler)
			}
		}

		// limited routes may be accessed using temporary jwt tokens
		v1Limited := privateGroup.Group("/v1/limited", mwTmpAuth)
		{
			authGroup := v1Limited.Group("/auth")
			{
				// GET /users/private/v1/limited/auth/me
				authGroup.GET("/me", authHandler.MeLimitedHandler)
			}

			usersGroup := v1Limited.Group("/users", middlewares.SetUIDParam, mwRequestedUser)
			{
				// PUT /users/private/v1/limited/users/profile
				usersGroup.PUT("/profile", usersHandler.UpdateHandler)
				// PUT /users/private/v1/limited/users/security-questions-answers
				usersGroup.GET("/security-questions-answers", securityQuestionsAnswersHandler.GetSecurityQuestionAnswersHandler)
				// GET /users/private/v1/limited/users/security-questions-answers
				usersGroup.PUT("/security-questions-answers", securityQuestionsAnswersHandler.UpdateSecurityQuestionAnswersHandler)
				// POST /users/private/v1/limited/users/generate-new-phone-code
				usersGroup.POST("/generate-new-phone-code", usersHandler.GenerateNewPhoneCode)
				// PUT /users/private/v1/limited/users/check-phone-code
				usersGroup.PUT("/check-phone-code", usersHandler.CheckPhoneCode)
				// POST /users/private/v1/limited/users/generate-new-email-code
				usersGroup.POST("/generate-new-email-code", usersHandler.GenerateNewEmailCode)
				// PUT /users/private/v1/limited/users/check-email-code
				usersGroup.PUT("/check-email-code", usersHandler.CheckEmailCode)
			}

			verificationsGroup := v1Limited.Group("/verifications")
			{
				verificationsGroup.GET("", verificationsHandler.ListWithTmpHandler)
				verificationsGroup.POST("", verificationsHandler.CreateHandler)
			}
		}
	}

	/*
	 |---------------------------------------------------
	 | Public router group
	 |---------------------------------------------------
	*/
	publicGroup := apiGroup.Group("/public")
	{
		v1Group := publicGroup.Group("/v1")
		{
			authGroup := v1Group.Group("auth")
			{
				// POST /users/public/v1/auth/signup
				authGroup.POST("/signup", authHandler.SimpleSignUpHandler)
				// POST /users/public/v1/auth/signin
				authGroup.POST("/signin", authHandler.SignInHandler)
				// POST /users/public/v1/auth/forgot-password
				authGroup.POST("/forgot-password", mwMaintenance, authHandler.ForgotPassword)
				// POST /users/public/v1/auth/reset-password
				authGroup.POST("/reset-password", mwMaintenance, authHandler.ResetPassword)
				// GET /users/public/v1/auth/refresh
				authGroup.GET("/refresh", mwMaintenance, authHandler.RefreshHandler)

				// mwUserFromCodeDisposable gives one time access to the given route using confirmation code
				mwUserFromCodeDisposable := middlewares.UserFromConfirmationCode(
					responseService,
					confirmationCodeService,
					models.ConfirmationCodeSubjectSetPassword,
					true,
					logger.New("middleware", "UserFromConfirmationCode"),
				)
				// POST /users/public/v1/auth/set-password
				authGroup.POST("/set-password", mwMaintenance, mwUserFromCodeDisposable, authHandler.SetPasswordHandler)

				// mwUserFromCode gives access to the given route using confirmation code
				mwUserFromCode := middlewares.UserFromConfirmationCode(
					responseService,
					confirmationCodeService,
					models.ConfirmationCodeSubjectSetPassword,
					false,
					logger.New("middleware", "UserFromConfirmationCode"),
				)
				// GET /users/public/v1/auth/confirmation-code/:code
				authGroup.GET("/confirmation-code/:code", mwMaintenance, mwUserFromCode, authHandler.GetConfirmationCodeHandler)
			}

			securityQuestionsGroup := v1Group.Group("security-questions", mwMaintenance)
			{
				// GET /users/public/v1/security-questions
				securityQuestionsGroup.GET("", securityQuestionHandler.ListHandler)
			}
		}
	}

	// If route not found returns StatusNotFound
	r.NoRoute(NotFound)

	// Handle OPTIONS request
	r.OPTIONS("/*cors", func(c *gin.Context) {
		c.Status(http.StatusOK)
		c.Abort()
		return
	})

	return r
}

// NotFound returns 404 NotFound
func NotFound(c *gin.Context) {
	e := responses.NewError().
		SetCode(responses.NotFound).
		SetTitleByCode(responses.NotFound).
		SetTarget(responses.TargetCommon).
		SetDetails("Not Found.")
	r := responses.NewResponse().SetStatusByCode(responses.NotFound).AddError(e)
	c.AbortWithStatusJSON(r.Status, r)
	return
}
