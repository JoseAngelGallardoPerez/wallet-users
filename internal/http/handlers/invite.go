package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services/invites"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/users"
	"github.com/Confialink/wallet-users/internal/validators"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
)

type InvitesHandler struct {
	repository           repositories.RepositoryInterface
	responseService      responses.ResponseHandler
	notificationsService *notifications.Notifications
	inviteCreator        *invites.Creator
	userCreator          *users.UserService
	logger               log15.Logger
}

func NewInvitesHandler(
	repository repositories.RepositoryInterface,
	responseService responses.ResponseHandler,
	notificationsService *notifications.Notifications,
	inviteCreator *invites.Creator,
	userCreator *users.UserService,
	logger log15.Logger,
) *InvitesHandler {
	return &InvitesHandler{
		repository,
		responseService,
		notificationsService,
		inviteCreator,
		userCreator,
		logger,
	}
}

func (h *InvitesHandler) CountHandler(ctx *gin.Context) {
	user, ok := ctx.Get("_user")
	if !ok {
		h.responseService.Forbidden(ctx)
		return
	}

	count, err := h.repository.GetInvitesRepository().CountByUserUID(user.(*userpb.User).UID)
	if err != nil {
		h.logger.Error("failed to count invites", "error", err)
		h.responseService.Error(ctx, responses.CannotRetrieveCollection, "Can't count invites")
		return
	}

	h.responseService.SuccessResponse(ctx, http.StatusOK, count)
	return
}

func (h *InvitesHandler) CreateHandler(ctx *gin.Context) {
	user, ok := ctx.Get("_user")
	if !ok {
		h.responseService.Forbidden(ctx)
		return
	}

	validator := validators.InviteValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		h.responseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	tx := h.repository.GetUsersRepository().DB.Begin()
	txUserRepo := h.repository.GetUsersRepository().WrapContext(tx)

	uid := user.(*userpb.User).UID
	email := validator.Data.Email
	password := validator.UserModel.Password

	currentUser, err := txUserRepo.FindByUID(uid)
	if err != nil {
		tx.Rollback()
		h.logger.Error("failed to find a user", "error", err)
		h.responseService.Error(ctx, responses.NotFound, "Can't find a user")
		return
	}

	validator.UserModel.Status = models.StatusActive
	validator.UserModel.SetCompanyDetails(currentUser.GetCompanyDetails())
	validator.UserModel.ClassId = currentUser.ClassId
	validator.UserModel.UserGroupId = currentUser.UserGroupId
	validator.UserModel.RoleName = currentUser.RoleName

	createdUser, err := h.userCreator.Create(&validator.UserModel, true, false, tx)
	if err != nil {
		tx.Rollback()
		h.logger.Error("failed to create a user", "error", err)
		h.responseService.Error(ctx, responses.CanNotCreateUser, "Can't create a user")
		return
	}

	invite, typedErr := h.inviteCreator.Call(email, uid, tx)
	if typedErr != nil {
		tx.Rollback()
		h.logger.Error("failed to create an invite", "error", typedErr)
		h.responseService.Error(ctx, responses.CanNotAddInvite, "Can't create an invite")
		return
	}

	currentUser.Status = models.StatusActive

	err = txUserRepo.UpdateStatusInfo(currentUser)
	if err != nil {
		tx.Rollback()
		h.logger.Error("failed to approve a user", "error", err)
		h.responseService.Error(ctx, responses.CanNotApproveUser, "Can't approve a user")
		return
	}

	tx.Commit()

	if _, err = h.notificationsService.InviteCreated(createdUser.UID, password); err != nil {
		h.logger.Error("failed to send notification", "error", err)
	}

	h.responseService.SuccessResponse(ctx, http.StatusCreated, invite)
	return
}
