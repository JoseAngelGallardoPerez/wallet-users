package middlewares

import (
	"net/http"

	"github.com/Confialink/wallet-users/internal/http/responses"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
)

func OwnerOrAdminOrRoot() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve current user
		user, ok := ctx.Get("_user")
		if !ok {
			// Returns a "403 StatusForbidden" response
			status := http.StatusForbidden
			ctx.Status(status)
			e := responses.NewError().SetTitleByStatus(status).SetDetails("User is not logged in.")
			r := responses.NewResponse().SetStatus(status).AddError(e)
			ctx.AbortWithStatusJSON(status, r)
			return
		}

		if ctx.Params.ByName("uid") != user.(*userpb.User).UID {
			rolename := user.(*userpb.User).RoleName
			if rolename != "admin" && rolename != "root" {
				// Returns a "403 StatusForbidden" response
				status := http.StatusForbidden
				ctx.Status(status)
				e := responses.NewError().SetTitleByStatus(status).SetDetails("User is not logged in.")
				r := responses.NewResponse().SetStatus(status).AddError(e)
				ctx.AbortWithStatusJSON(status, r)
				return
			}
		}
	}
}
