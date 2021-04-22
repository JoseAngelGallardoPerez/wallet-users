package forms

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ResetPassword struct {
	NewPassword     string `binding:"required,min=8,max=128,specialCharacterRequired,numberRequired,uppercaseLetterRequired,lowercaseLetterRequired"`
	ConfirmPassword string `binding:"required,eqfield=NewPassword"`
}

// BindJSON binding from JSON
func (s *ResetPassword) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}
	return nil
}
