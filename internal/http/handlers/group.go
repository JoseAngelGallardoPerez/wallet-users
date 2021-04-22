package handlers

import (
	"net/http"

	"github.com/Confialink/wallet-users/internal/validators"

	"github.com/gin-gonic/gin"

	"strconv"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
)

type UserGroupsService struct {
	Repository      repositories.RepositoryInterface
	ResponseService responses.ResponseHandler
}

// UserGroupsServiceHandler is interface for handler functionality
// that ought to be implemented manually.
type UserGroupsServiceHandler interface {
	ListHandler(ctx *gin.Context)
	GetHandler(ctx *gin.Context)
	CreateHandler(ctx *gin.Context)
	UpdateHandler(ctx *gin.Context)
}

func NewUserGroupsService(repository repositories.RepositoryInterface, responseService responses.ResponseHandler) *UserGroupsService {
	return &UserGroupsService{
		repository,
		responseService,
	}
}

// GetHandler returns user by uid
func (srv *UserGroupsService) GetHandler(ctx *gin.Context) {
	userGroup, err := srv.Repository.GetUserGroupsRepository().FindById(srv.getIdParam(ctx))
	if err != nil {
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	// Returns a "200 OK" response
	srv.ResponseService.OkResponse(ctx, userGroup)
	return
}

// ListHandler returns list of user groups
func (srv *UserGroupsService) ListHandler(ctx *gin.Context) {
	limitQuery := ctx.DefaultQuery("limit", "10")
	pageQuery := ctx.DefaultQuery("page", "1")

	query := srv.Repository.GetUserGroupsRepository().Filter(ctx.Request.URL.Query())

	pagination, err := srv.Repository.GetUserGroupsRepository().Paginate(query, pageQuery, limitQuery)
	if err != nil {
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of user groups")
		return
	}

	// Returns a "200 OK" response
	srv.ResponseService.OkResponse(ctx, pagination)
}

// CreateHandler creates new user group
func (srv *UserGroupsService) CreateHandler(ctx *gin.Context) {
	// Checks if the query entry is valid
	validator := validators.CreateUserGroupValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	userGroupModel := validator.UserGroupModel

	// Check the database if there is a group with such name.
	record, err := srv.Repository.GetUserGroupsRepository().FindByName(userGroupModel.Name)
	if record != nil {
		// Returns a "409 StatusConflict" response
		details := "User group with provided name already exists in database."
		srv.ResponseService.Error(ctx, responses.UserGroupAlreadyExists, details)
		return
	}

	// Create new user group in DB
	err = srv.Repository.GetUserGroupsRepository().Create(&userGroupModel)
	if err != nil {
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CanNotCreateUserGroup, "Can't create a user group")
		return
	}

	// Returns a "201 Created" response
	srv.ResponseService.SuccessResponse(ctx, http.StatusCreated, userGroupModel)
}

// UpdateHandler updates a user group
func (srv *UserGroupsService) UpdateHandler(ctx *gin.Context) {
	userGroup, err := srv.Repository.GetUserGroupsRepository().FindById(srv.getIdParam(ctx))
	if err != nil {
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	// Checks if the query entry is valid
	validator := validators.GetUpdateUserGroupValidatorFillWith(userGroup)
	if err = validator.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	// Update a user group
	err = srv.Repository.GetUserGroupsRepository().Update(userGroup, &validator.UserGroupModel)
	if err != nil {
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CanNotUpdateUserGroup, "Can't update a user group")
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

// DeleteHandler deletes a user group
func (srv *UserGroupsService) DeleteHandler(ctx *gin.Context) {
	userGroup, err := srv.Repository.GetUserGroupsRepository().FindById(srv.getIdParam(ctx))
	if err != nil {
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	count, err := srv.Repository.GetUsersRepository().CountByUserGroupID(userGroup.ID)
	if err != nil {
		// Returns a "500 StatusInternalServerError" response
		details := "Can't count users of a user group"
		srv.ResponseService.Error(ctx, responses.CanNotDeleteUserGroup, details)
		return
	}

	if count > 0 {
		// Returns a "500 StatusInternalServerError" response
		details := "Can't delete a user group, user group contains users"
		srv.ResponseService.Error(ctx, responses.CanNotDeleteUserGroup, details)
		return
	}

	// Update a user group
	err = srv.Repository.GetUserGroupsRepository().Delete(userGroup)
	if err != nil {
		// Returns a "500 StatusInternalServerError" response
		srv.ResponseService.Error(ctx, responses.CanNotDeleteUserGroup, "Can't delete a user group")
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}

// getIdParam returns id or nil
func (srv UserGroupsService) getIdParam(ctx *gin.Context) uint64 {
	id := ctx.Params.ByName("id")

	// convert string to uint
	id64, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		srv.ResponseService.Error(ctx, responses.InvalidIdForUserGroup, "Id param must be an integer")
		return 0
	}

	return uint64(id64)
}
