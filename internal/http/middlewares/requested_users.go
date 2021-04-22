package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
)

func RequestedUser(repo *repositories.UsersRepository, responseService responses.ResponseHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Params.ByName("uid")

		foundUser, err := repo.FindByUID(uid)
		if err != nil {
			responseService.NotFound(ctx)
			ctx.Abort()
			return
		}
		ctx.Set("_requested_user", foundUser)
	}
}
