package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// UnblockUser is
type UnblockUser struct {
	UID string `json:"uid" binding:"required"`
}

// UnblockUserValidator is
type UnblockUsersValidator struct {
	Data []*UnblockUser `json:"data" binding:"required,dive,required"`
}

// BindJSON binding from JSON
func (s *UnblockUsersValidator) BindJSON(c *gin.Context) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	err := c.ShouldBindWith(s, b)
	if err != nil {
		return err
	}
	return nil
}
