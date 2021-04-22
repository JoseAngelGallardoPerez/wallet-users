package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// UnblockIP is
type UnblockIP struct {
	ID uint64 `json:"id" binding:"required"`
}

// UnblockIpsValidator is
type UnblockIpsValidator struct {
	Data []*UnblockIP `json:"data" binding:"required,dive,required"`
}

// BindJSON binding from JSON
func (s *UnblockIpsValidator) BindJSON(c *gin.Context) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	err := c.ShouldBindWith(s, b)
	if err != nil {
		return err
	}
	return nil
}
