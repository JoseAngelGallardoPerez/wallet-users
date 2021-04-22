package middlewares

import (
	"net/http"

	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	httpAuth "github.com/Confialink/wallet-users/internal/http/services/auth"
	"github.com/Confialink/wallet-users/internal/services/auth"
)

func TmpAuthByToken(
	service *auth.TemporaryTokens,
	usersRepo *repositories.UsersRepository,
	logger log15.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		signedToken := c.Request.Header.Get(httpAuth.TmpAuthHeader)

		if signedToken == "" {
			c.Abort()
			c.Status(http.StatusUnauthorized)
			return
		}

		token, err := service.Verify(signedToken)
		if err != nil {
			logger.Error("token verification failed", "error", err)
			c.Abort()
			c.Status(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		user, err := usersRepo.FindByUID(claims["uid"].(string))

		if err != nil {
			logger.Error("failed to retrieve user", "error", err)
			c.Abort()
			c.Status(http.StatusUnauthorized)
		}
		groupId := uint64(0)
		if user.UserGroupId != nil {
			groupId = *user.UserGroupId
		}
		pbuser := &users.User{
			UID:         user.UID,
			Email:       user.Email,
			Username:    user.Username,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			RoleName:    user.RoleName,
			GroupId:     groupId,
			CompanyName: user.CompanyDetails.CompanyName,
			PhoneNumber: user.PhoneNumber,
		}

		c.Set("_user", pbuser)
	}
}
