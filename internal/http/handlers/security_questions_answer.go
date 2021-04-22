package handlers

import (
	"fmt"
	"net/http"

	"github.com/Confialink/wallet-users/internal/validators"

	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/gin-gonic/gin"
)

// SettingService is settings service for CRUD operations.
type SecurityQuestionAnswerService struct {
	Repository      repositories.RepositoryInterface
	ResponseService responses.ResponseHandler
	Db              *gorm.DB
}

func NewSecurityQuestionAnswerService(
	repository repositories.RepositoryInterface,
	responseService responses.ResponseHandler,
	db *gorm.DB,
) *SecurityQuestionAnswerService {
	return &SecurityQuestionAnswerService{
		repository,
		responseService,
		db,
	}
}

// GetSecurityQuestionAnswersHandler returns list of security question answers
func (s SecurityQuestionAnswerService) GetSecurityQuestionAnswersHandler(c *gin.Context) {
	user := GetRequestedUser(c)
	if user == nil {
		s.ResponseService.NotFound(c)
		return
	}

	answers, _ := s.Repository.GetSecurityQuestionsAnswerRepository().FindByUID(user.UID)
	// Returns a "200 OK" response
	s.ResponseService.OkResponse(c, answers)
}

// UpdateSecurityQuestionAnswersHandler updates security question answers
func (s SecurityQuestionAnswerService) UpdateSecurityQuestionAnswersHandler(c *gin.Context) {
	user := GetRequestedUser(c)
	if user == nil {
		s.ResponseService.NotFound(c)
		return
	}
	var errs []*responses.Error

	// Checks if the query entry is valid
	validator := validators.UpdateSecurityQuestionsAnswersValidator{}
	if err := validator.BindJSON(c); err != nil {
		// Returns a "422 StatusUnprocessableEntity" response
		s.ResponseService.ValidatorErrorResponse(c, responses.UnprocessableEntity, err)
		return
	}

	tx := s.Db.Begin()
	repo := s.Repository.GetSecurityQuestionsAnswerRepository().WrapContext(tx)
	// at some moment users had many security questions-answers so we have to remove all extra answers
	// now users have only one security answer
	repo.DeleteByUID(user.UID)

	// Update answers
	for _, item := range validator.Data {
		item.UID = user.UID
		_, err := repo.FirstOrCreate(item)
		if err != nil {
			tx.Rollback()
			// Returns a "500 StatusInternalServerError" response
			e := responses.NewCommonError().
				ApplyCode(responses.CannotUpdateSecurityQuestionAnswer).
				SetDetails(fmt.Sprintf("Can't update security question answer `%d`", item.AID))
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		tx.Rollback()
		s.ResponseService.Errors(c, http.StatusInternalServerError, errs)
		return
	}

	tx.Commit()
	// Returns a "200 OK" response
	s.ResponseService.OkResponse(c, validator.Data)
}
