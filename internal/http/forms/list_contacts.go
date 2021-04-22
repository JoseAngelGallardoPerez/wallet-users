package forms

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ListContacts struct {
	PhoneNumbers []string `json:"phoneNumbers" binding:"required"`
}

// BindJSON binding from JSON
func (s *ListContacts) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}
	return nil
}
