package handlers

import (
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/gin-gonic/gin"
)

// SecurityQuestionService is security question service for CRUD operations.
type SecurityQuestionService struct {
	Repository      repositories.RepositoryInterface
	ResponseService responses.ResponseHandler
}

func NewSecurityQuestionService(repository repositories.RepositoryInterface, responseService responses.ResponseHandler) *SecurityQuestionService {
	return &SecurityQuestionService{
		repository,
		responseService,
	}
}

// ListHandler returns the list of security question
func (srv *SecurityQuestionService) ListHandler(ctx *gin.Context) {
	questions, err := srv.Repository.GetSecurityQuestionRepository().GetAll()
	if err != nil {
		// Returns a "404 StatusNotFound" response
		srv.ResponseService.NotFound(ctx)
		return
	}

	// Returns a "200 OK" response
	srv.ResponseService.OkResponse(ctx, questions)
	return
}
