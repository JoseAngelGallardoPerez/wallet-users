package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LoginDeviceValidator struct {
	Pin      string `json:"pin" binding:"required,min=5,max=5"`
	DeviceId string `json:"deviceId" binding:"required,max=255"`
}

// BindJSON binding from JSON
func (s *LoginDeviceValidator) BindJSON(ctx *gin.Context) error {
	b := binding.Default(ctx.Request.Method, ctx.ContentType())

	err := ctx.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	return nil
}
