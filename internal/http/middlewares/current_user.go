package middlewares

import (
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/users"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

func UserFromAccessToken(responseService responses.ResponseHandler, usersRepository *repositories.UsersRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, ok := c.Get("AccessToken")
		if !ok {
			// Returns a "401 StatusUnauthorized" response
			responseService.Error(c, responses.Unauthorized, "Access token not found.")
			return
		}

		user, err := usersRepository.FindUserByTokenAndSubject(accessToken.(string), auth.ClaimAccessSub)
		if err != nil {
			// Returns a "400 StatusBadRequest" response
			responseService.Error(c, responses.CannotFindUserByAccessToken, "Can't find user.")
			return
		}
		c.Set("_current_user", user)
	}
}

// UserFromConfirmationCode retrieves user from confirmation code and set it to context, removes the code if oneTime is set to true
func UserFromConfirmationCode(
	responseService responses.ResponseHandler,
	confirmationCodeService *users.ConfirmationCode,
	confirmationCodeSubject string,
	oneTime bool,
	logger log15.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			code = c.Params.ByName("code")
			if code == "" {
				responseService.Error(c, responses.ConfirmationCodeIsInvalid, "Code is not found.")
				return
			}
		}

		confirmationCode, err := confirmationCodeService.FindByCodeAndSubject(code, confirmationCodeSubject)
		if err != nil {
			responseService.Error(c, responses.ConfirmationCodeIsInvalid, "Code is invalid or already used.")
			return
		}

		if confirmationCode.ExpiresAt.Before(time.Now()) {
			responseService.Error(c, responses.ConfirmationCodeIsInvalid, "Code is expired.")
			if err := confirmationCodeService.DeleteConfirmationCode(confirmationCode); err != nil {
				logger.Error("cannot remove expired confirmation code", "err", err)
			}
			return
		}

		c.Set("_current_user", confirmationCode.User)
		c.Set("_confirmation_code", confirmationCode)

		if oneTime {
			if err := confirmationCodeService.DeleteConfirmationCode(confirmationCode); err != nil {
				logger.Error("cannot remove confirmation code", "err", err)
			}
		}
	}
}

func SetUIDParam(c *gin.Context) {
	uid := ""

	if user, ok := c.Get("_user"); ok {
		if pbuser, ok := user.(*pb.User); ok {
			uid = pbuser.UID
		} else {
			panic("unexpected _user type")
		}
	}

	if user, ok := c.Get("_current_user"); ok {
		if currentUser, ok := user.(*models.User); ok {
			uid = currentUser.UID
		} else {
			panic("unexpected _current_user type")
		}
	}

	if uid != "" {
		c.Params = append(c.Params, gin.Param{Key: "uid", Value: uid})
	}
}
