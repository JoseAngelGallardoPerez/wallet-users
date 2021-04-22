package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ChangePassword struct {
	PreviousPassword string
	ProposedPassword string
	ConfirmPassword  string
}

// ChangePasswordValidator is validator for change password
type ChangePasswordValidator struct {
	Data struct {
		PreviousPassword string `json:"previousPassword" binding:"required"`
		ProposedPassword string `json:"proposedPassword" binding:"required,min=8,max=255,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`
		ConfirmPassword  string `json:"confirmPassword" binding:"required,eqfield=ProposedPassword"`
	} `json:"data"`
	Model ChangePassword `json:"-"`
}

type SetPassword struct {
	ProposedPassword string `json:"proposedPassword" binding:"required,min=8,max=255,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`
	ConfirmPassword  string `json:"confirmPassword" binding:"required,eqfield=ProposedPassword"`
}

// ResetPassword struct contains validation rule for reset password
type ResetPassword struct {
	ConfirmationCode string `json:"confirmationCode" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required,min=8,max=255,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`
}

// BindJSON binding from JSON
func (s *ChangePasswordValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	s.Model.PreviousPassword = s.Data.PreviousPassword
	s.Model.ProposedPassword = s.Data.ProposedPassword

	return nil
}
