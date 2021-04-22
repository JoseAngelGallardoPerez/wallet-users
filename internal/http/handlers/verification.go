package handlers

import (
	"net/http"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/services/verification"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	verificationStates "github.com/Confialink/wallet-users/internal/services/verification_states"
)

type VerificationHandler struct {
	repository      repositories.RepositoryInterface
	responseService responses.ResponseHandler
	creator         *verification.Creator
	validator       *verification.Validator
	logger          log15.Logger
}

func NewVerificationHandler(
	repository repositories.RepositoryInterface,
	responseService responses.ResponseHandler,
	creator *verification.Creator,
	validator *verification.Validator,
	logger log15.Logger,
) *VerificationHandler {
	return &VerificationHandler{
		repository,
		responseService,
		creator,
		validator,
		logger,
	}
}

func (h *VerificationHandler) ListHandler(ctx *gin.Context) {
	uid, ok := ctx.Params.Get("uid")
	if !ok {
		h.responseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of verifications")
		return
	}

	verifications, err := h.repository.GetVerificationRepository().FindByUID(uid)
	if err != nil {
		h.logger.Error("failed to load verifications", "error", err)
		h.responseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of verifications")
		return
	}

	h.responseService.OkResponse(ctx, verifications)
}

func (h *VerificationHandler) ListWithTmpHandler(ctx *gin.Context) {
	user, ok := ctx.Get("_user")
	if !ok {
		h.responseService.Forbidden(ctx)
		return
	}

	verifications, err := h.repository.GetVerificationRepository().FindByUID(user.(*userpb.User).GetUID())
	if err != nil {
		h.logger.Error("failed to load verifications", "error", err)
		h.responseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of verifications")
		return
	}

	h.responseService.OkResponse(ctx, verifications)
}

func (h *VerificationHandler) RequestHandler(ctx *gin.Context) {
	user, ok := ctx.Get("_user")
	if !ok {
		h.responseService.Forbidden(ctx)
		return
	}

	verification, err := h.getVerification(ctx)
	if err != nil {
		h.responseService.Error(ctx, responses.VerificationNotFound, "Can't find a verification")
		return
	}

	if verification.UserUID != user.(*userpb.User).GetUID() {
		h.responseService.Forbidden(ctx)
		return
	}

	if err := verificationStates.NewVerificationState(verification).HandleVerificationRequest(); err != nil {
		h.responseService.Error(ctx, responses.CanNotCreateVerificationRequest, err.Error())
		return
	}

	err = h.repository.GetVerificationRepository().Save(verification)
	if err != nil {
		h.logger.Error("failed to update verification", "error", err)
		h.responseService.Error(ctx, responses.CanNotUpdateVerification, "Can't update a verification")
		return
	}

	h.responseService.SuccessResponse(ctx, http.StatusOK, verification)
	return
}

func (h *VerificationHandler) ApproveHandler(ctx *gin.Context) {
	verification, err := h.getVerification(ctx)
	if err != nil {
		h.responseService.Error(ctx, responses.VerificationNotFound, "Can't find a verification.")
		return
	}

	if err := verificationStates.NewVerificationState(verification).HandleAdminApprove(); err != nil {
		h.responseService.Error(ctx, responses.CanNotApproveVerificationRequest, err.Error())
		return
	}

	err = h.repository.GetVerificationRepository().Save(verification)
	if err != nil {
		h.logger.Error("failed to update verification", "error", err)
		h.responseService.Error(ctx, responses.CanNotUpdateVerification, "Can't update a verification")
		return
	}

	h.responseService.SuccessResponse(ctx, http.StatusOK, verification)
	return
}

func (h *VerificationHandler) CancelHandler(ctx *gin.Context) {
	verification, err := h.getVerification(ctx)
	if err != nil {
		h.responseService.Error(ctx, responses.VerificationNotFound, "Can't find a verification.")
		return
	}

	if err := verificationStates.NewVerificationState(verification).HandleAdminCancellation(); err != nil {
		h.responseService.Error(ctx, responses.CanNotCancelVerificationRequest, err.Error())
		return
	}

	err = h.repository.GetVerificationRepository().Save(verification)
	if err != nil {
		h.logger.Error("failed to update verification", "error", err)
		h.responseService.Error(ctx, responses.CanNotUpdateVerification, "Can't update a verification")
		return
	}

	h.responseService.SuccessResponse(ctx, http.StatusOK, verification)
	return
}

func (h *VerificationHandler) CreateHandler(ctx *gin.Context) {
	user, ok := ctx.Get("_user")
	if !ok {
		h.responseService.Forbidden(ctx)
		return
	}

	form := struct {
		VerificationType string  `json:"verificationType" binding:"required,verificationTypeOneOf"`
		FileId           *uint64 `json:"fileId" binding:"required"`
	}{}

	if err := ctx.ShouldBindJSON(&form); err != nil {
		errors.AddShouldBindError(ctx, err)
		h.responseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	uid := user.(*userpb.User).UID

	if typedErr := h.validator.CanCreate(form.VerificationType, uid); typedErr != nil {
		errors.AddErrors(ctx, typedErr)
		return
	}

	model, err := h.creator.Call(form.VerificationType, uid, *form.FileId)
	if err != nil {
		h.logger.Error("failed to create verification", "error", err)
		h.responseService.Error(ctx, responses.CanNotAddVerification, "Can't create a verification")
		return
	}

	h.responseService.SuccessResponse(ctx, http.StatusCreated, model)
	return
}

func (h *VerificationHandler) getVerification(ctx *gin.Context) (*models.Verification, error) {
	id, err := getUint64Param(ctx, "id")
	if err != nil {
		return nil, err
	}

	verificationModel, err := h.repository.GetVerificationRepository().FindById(id)
	if err != nil {
		h.logger.Error("unable to retrieve verification", "error", err)
		return nil, err
	}

	return verificationModel, nil
}
