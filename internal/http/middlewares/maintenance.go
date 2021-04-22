package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-users/internal/services/syssettings"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
)

func Maintenance(sysSettings *syssettings.SysSettings) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Get("_user")
		if ok {
			r := user.(*userpb.User).RoleName
			if r == "root" {
				return
			}
		}

		maintenanceModeSettings, err := sysSettings.GetMaintenanceModeSettings()
		if err != nil {
			ctx.Status(http.StatusForbidden)
			ctx.Abort()
			return
		}

		if maintenanceModeSettings.Enabled {
			ctx.Status(http.StatusForbidden)
			ctx.Abort()
			return
		}
	}
}
