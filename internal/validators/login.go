package validators

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// LoginValidator is validator for login
type LoginValidator struct {
	Data struct {
		Email    string `json:"email" binding:"required,min=3,max=255"`
		Password string `json:"password" binding:"required"`
	} `json:"data"`
	UserModel models.User `json:"-"`
}

// BindJSON binding from JSON
func (s *LoginValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.UserModel.Email = s.Data.Email
	s.UserModel.Password = s.Data.Password

	return nil
}
