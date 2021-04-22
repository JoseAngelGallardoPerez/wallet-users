package handlers

import (
	"fmt"
	"net/http"

	"github.com/Confialink/wallet-users/internal/validators"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/gin-gonic/gin"
)

// BlockedIpsService is blocked ips service for CRUD operations.
type BlockedIpsService struct {
	Repository      repositories.RepositoryInterface
	ResponseService responses.ResponseHandler
}

func NewBlockedIpsService(repository repositories.RepositoryInterface, responseService responses.ResponseHandler) *BlockedIpsService {
	return &BlockedIpsService{
		Repository:      repository,
		ResponseService: responseService,
	}
}

// ListHandler returns list of blocked ips
func (srv *BlockedIpsService) ListHandler(ctx *gin.Context) {
	limitQuery := ctx.DefaultQuery("limit", "10")
	pageQuery := ctx.DefaultQuery("page", "1")

	query := srv.Repository.GetBlockedIpsRepository().Filter(ctx.Request.URL.Query())

	pagination, err := srv.Repository.GetBlockedIpsRepository().Paginate(query, pageQuery, limitQuery)
	if err != nil {
		// Returns a "400 StatusBadRequest" response
		srv.ResponseService.Error(ctx, responses.CannotRetrieveCollection, "Can't load list of blocked ips")
		return
	}

	// Returns a "200 OK" response
	srv.ResponseService.OkResponse(ctx, pagination)
}

// UnblockIpsHandler unblock selected ips
func (srv *BlockedIpsService) UnblockIpsHandler(ctx *gin.Context) {
	var errs []*responses.Error

	// Checks if the query entry is valid
	validator := validators.UnblockIpsValidator{}
	if err := validator.BindJSON(ctx); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		srv.ResponseService.ValidatorErrorResponse(ctx, responses.UnprocessableEntity, err)
		return
	}

	// Unblock ips
	for _, item := range validator.Data {
		err := srv.Repository.GetBlockedIpsRepository().Delete(&models.BlockedIp{ID: item.ID})
		if err != nil {
			e := responses.NewCommonError().
				ApplyCode(responses.CannotUnblockIp).
				SetDetails(fmt.Sprintf("Can't unblock ip `%d`", item.ID))
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		srv.ResponseService.Errors(ctx, http.StatusInternalServerError, errs)
		return
	}

	// Returns a "204 StatusNoContent" response
	ctx.JSON(http.StatusNoContent, nil)
}
