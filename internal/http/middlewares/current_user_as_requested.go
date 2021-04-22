package middlewares

import (
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
)

// CurrentUserAsRequested finds a user in DB and adds him to the request context
func CurrentUserAsRequested(repo *repositories.UsersRepository, responseService responses.ResponseHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Get("_user")
		if !ok {
			// Returns a "400 StatusBadRequest" response
			responseService.Error(ctx, responses.CannotFindUserByAccessToken, "Can't find user.")
			return
		}

		foundUser, err := repo.FindByUID(user.(*userpb.User).UID)
		if err != nil {
			responseService.NotFound(ctx)
			ctx.Abort()
			return
		}
		ctx.Set("_requested_user", foundUser)
	}
}
