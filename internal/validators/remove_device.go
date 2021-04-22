package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RemoveDeviceValidator struct {
	DeviceID string `json:"deviceId" binding:"required,max=255"`
}

// BindJSON binding from JSON
func (s *RemoveDeviceValidator) BindJSON(c *gin.Context) error {
	b := binding.Default(c.Request.Method, c.ContentType())

	err := c.ShouldBindWith(s, b)
	if err != nil {
		return err
	}

	return nil
}
