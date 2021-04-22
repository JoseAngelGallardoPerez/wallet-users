package validators

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// UpdateUserValidator ...
type UpdateUserGroupValidator struct {
	Data struct {
		Name        string `json:"name" binding:"required,min=4,max=255,noGroupNameExists"`
		Description string `json:"description,omitempty"`
	} `json:"data"`
	UserGroupModel models.UserGroup `json:"-"`
}

// GetUpdateUserValidator ...
func GetUpdateUserGroupValidator() UpdateUserGroupValidator {
	return UpdateUserGroupValidator{}
}

// GetUpdateUserValidatorFillWith ...
func GetUpdateUserGroupValidatorFillWith(userGroup *models.UserGroup) UpdateUserGroupValidator {

	v := GetUpdateUserGroupValidator()
	v.Data.Name = userGroup.Name
	v.Data.Description = userGroup.Description
	v.UserGroupModel.ID = userGroup.ID
	v.UserGroupModel.Name = userGroup.Name
	v.UserGroupModel.Description = userGroup.Description

	return v
}

// BindJSON binding from JSON
func (s *UpdateUserGroupValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.UserGroupModel.Name = s.Data.Name
	s.UserGroupModel.Description = s.Data.Description

	return nil
}
