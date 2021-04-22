package validators

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// CreateUserGrouValidator is
type CreateUserGroupValidator struct {
	Data struct {
		Name        string `json:"name" binding:"required,min=4,max=255"`
		Description string `json:"description,omitempty"`
	} `json:"data"`
	UserGroupModel models.UserGroup `json:"-"`
}

// BindJSON binding from JSON
func (s *CreateUserGroupValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.UserGroupModel.Name = s.Data.Name
	s.UserGroupModel.Description = s.Data.Description

	return nil
}
