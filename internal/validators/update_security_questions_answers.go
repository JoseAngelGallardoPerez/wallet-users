package validators

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// UpdateSecurityQuestionsAnswersValidator is
type UpdateSecurityQuestionsAnswersValidator struct {
	Data []*models.SecurityQuestionsAnswer `json:"data" binding:"required,dive,required"`
}

// BindJSON binding from JSON
func (s *UpdateSecurityQuestionsAnswersValidator) BindJSON(c *gin.Context) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	err := c.ShouldBindWith(s, b)
	if err != nil {
		return err
	}
	return nil
}
