package validators

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ForgotPasswordValidator is validator for forgot password
type ForgotPasswordValidator struct {
	Data struct {
		Email string `json:"email" binding:"required"`
	} `json:"data"`
	UserModel models.User `json:"-"`
}

// BindJSON binding from JSON
func (s *ForgotPasswordValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.UserModel.Email = s.Data.Email

	return nil
}
