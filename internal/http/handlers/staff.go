package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/http/serializers"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/permissions"
	system_logs "github.com/Confialink/wallet-users/internal/services/system-logs"
	"github.com/Confialink/wallet-users/internal/services/users"
	"github.com/Confialink/wallet-users/internal/validators"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
)

// StaffsService contains handler for staff users API
type StaffsService struct {
	Repository           repositories.RepositoryInterface
	ResponseService      responses.ResponseHandler
	SystemLogsService    *system_logs.SystemLogsService
	notificationsService *notifications.Notifications
	PermissionsService   *permissions.Permissions
	userCreator          *users.UserService
	params               *HandlerParams
	logger               log15.Logger
	PasswordService      *services.Password
}

// NewStaffsService return new StaffsService
func NewStaffsService(
	repository repositories.RepositoryInterface,
	responseService responses.ResponseHandler,
	systemLogsService *system_logs.SystemLogsService,
	notificationsService *notifications.Notifications,
	permissionsService *permissions.Permissions,
	userCreator *users.UserService,
	params *HandlerParams,
	logger log15.Logger,
	PasswordService *services.Password,
) *StaffsService {
	return &StaffsService{
		repository,
		responseService,
		systemLogsService,
		notificationsService,
		permissionsService,
		userCreator,
		params,
		logger,
		PasswordService,
	}
}

func (srv *StaffsService) isMainCompany(currentUser *userpb.User, requestedUser *models.User) bool {
	return strings.EqualFold(currentUser.CompanyName, requestedUser.CompanyDetails.CompanyName)
}

// GetHandler returns staff user by uid
func (srv *StaffsService) GetHandler(ctx *gin.Context) {
	currentUser := GetCurrentUser(ctx)
	requestedUser := GetRequestedUser(ctx)

	isMainCompany := srv.isMainCompany(currentUser, requestedUser)
	if !isMainCompany {
		srv.ResponseService.Forbidden(ctx)
		return
	}

	serialized := serializers.NewGetUser(requestedUser).Serialize()
	srv.ResponseService.OkResponse(ctx, serialized)
	return
}

// ListHandler returns list of staff users
func (srv *StaffsService) ListHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "ListHandler")

	currentUser := GetCurrentUser(ctx)

	if hasPerm := srv.PermissionsService.CanViewUserProfile(currentUser.UID); !hasPerm {
		srv.ResponseService.Forbidden(ctx)
		return
	}

	limitQuery := ctx.DefaultQuery("limit", "10")
	pageQuery := ctx.DefaultQuery("page", "1")
	params := ctx.Request.URL.Query()

	query := srv.Repository.GetUsersRepository().Filter(params)
	query = srv.Repository.GetUsersRepository().FilterByCompanyName(query, currentUser.CompanyName)
	query = srv.Repository.GetUsersRepository().FilterByRoleName(query, currentUser.RoleName)

	pagination, err := srv.Repository.GetUsersRepository().Paginate(query, pageQuery, limitQuery, serializers.NewUsers())
	if err != nil {
		logger.Error("failed to load list of users", "error", err)
		srv.ResponseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of users")
		return
	}

	srv.ResponseService.OkResponse(ctx, pagination)
}

// CreateHandler creates a new staff user
func (srv *StaffsService) CreateHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "CreateHandler")

	user := GetCurrentUser(ctx)

	validator := validators.CreateStaffValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	uid := user.UID
	password := validator.UserModel.Password

	if hasPerm := srv.PermissionsService.CanCreateProfile(uid, &validator.UserModel); !hasPerm {
		srv.ResponseService.Forbidden(ctx)
		return
	}

	tx := srv.Repository.GetUsersRepository().DB.Begin()
	txUserRepo := srv.Repository.GetUsersRepository().WrapContext(tx)

	currentUser, err := txUserRepo.FindByUID(uid)
	if err != nil {
		tx.Rollback()
		logger.Error("failed to find a user", "error", err)
		srv.ResponseService.Error(ctx, responses.NotFound, "Can't find a user")
		return
	}

	validator.UserModel.ParentId = currentUser.UID
	validator.UserModel.Status = models.StatusActive
	validator.UserModel.SetCompanyDetails(currentUser.GetCompanyDetails())
	validator.UserModel.UserGroupId = currentUser.UserGroupId
	validator.UserModel.RoleName = currentUser.RoleName

	company, err := srv.Repository.GetCompanyRepository().FindByID(*currentUser.CompanyID)
	if err != nil {
		logger.Error("сan't get user company", "error", err)
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CanNotCreateUser, "сan't get user company")
		return
	}

	validator.UserModel.CompanyID = &company.ID
	validator.UserModel.CompanyDetails.ID = company.ID
	validator.UserModel.CompanyDetails.CompanyName = company.CompanyName
	validator.UserModel.CompanyDetails.DirectorLastName = company.DirectorLastName
	validator.UserModel.CompanyDetails.DirectorFirstName = company.DirectorFirstName
	validator.UserModel.CompanyDetails.CompanyRole = company.CompanyRole
	validator.UserModel.CompanyDetails.CompanyType = company.CompanyType

	// Create new user
	createdUser, err := srv.userCreator.Create(&validator.UserModel, true, false, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("failed to create a user", "error", err)
		srv.ResponseService.Error(ctx, responses.CanNotCreateUser, "Can't create a user")
		return
	}

	tx.Commit()

	if currentUser != nil {
		srv.SystemLogsService.LogCreateUserProfileAsync(createdUser, currentUser.UID)
	}

	if _, err = srv.notificationsService.ProfileCreated(createdUser.UID, password, ""); err != nil {
		logger.Error("failed to send notification", "error", err)
		// do not return an error
	}

	srv.ResponseService.SuccessResponse(ctx, http.StatusCreated, createdUser)
	return
}

// UpdateHandler updates a staff user
func (srv *StaffsService) UpdateHandler(ctx *gin.Context) {
	logger := srv.logger.New("action", "UpdateHandler")

	currentUser := GetCurrentUser(ctx)
	requestedUser := GetRequestedUser(ctx)

	isMainCompany := srv.isMainCompany(currentUser, requestedUser)
	if !isMainCompany {
		srv.ResponseService.Forbidden(ctx)
		return
	}

	validator := validators.GetUpdateStaffValidatorFillWith(*requestedUser)
	if err := validator.BindJSON(ctx); err != nil {
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	tx := srv.Repository.GetUsersRepository().DB.Begin()
	txUserRepo := srv.Repository.GetUsersRepository().WrapContext(tx)

	_, err := txUserRepo.Save(&validator.UserModel)
	if err != nil {
		tx.Rollback()
		logger.Error("failed to update user", "error", err)
		srv.ResponseService.Error(ctx, responses.CanNotUpdateUser, "Can't update a user")
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusNoContent, nil)
}
