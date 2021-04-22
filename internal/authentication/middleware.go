package authentication

import (
	"context"
	"log"
	"net/http"
	"strings"

	errorsPkg "github.com/Confialink/wallet-pkg-errors"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/srvdiscovery"
)

// Middleware authentication middleware
func Middleware(logger log15.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		accessToken, ok := ExtractToken(c)
		if !ok {
			c.Header("Authentication", `Bearer realm="private"`)
			logger.Info("Access token not found")
			_ = c.Error(accessTokenNotFoundError())
			c.Abort()
			return
		}

		client := userpb.NewUserHandlerProtobufClient(getRPCUsersServerAddr(), &http.Client{})
		res, err := client.ValidateAccessToken(context.Background(), &userpb.Request{AccessToken: accessToken})
		if nil != err {
			c.Header("Authentication", `Bearer realm="private"`)
			logger.Info("Access token invalid", "err", err)
			_ = c.Error(accessTokenInvalidError())
			c.Abort()
			return
		}

		c.Set("AccessToken", accessToken)
		c.Set("_user", res.User)
	}
}

// ExtractToken extracts jwt token from the header "Authorization" field with Bearer
func ExtractToken(c *gin.Context) (string, bool) {
	tokens := c.Request.Header.Get("Authorization")
	if len(tokens) < 8 || !strings.EqualFold(tokens[0:7], "Bearer ") {
		return "", false // empty token
	}

	return tokens[7:], true
}

func getRPCUsersServerAddr() string {
	endpoint, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameUsers)
	if nil != err {
		log.Fatalf(err.Error())
	}
	return endpoint.String()
}

func accessTokenNotFoundError() *errorsPkg.PublicError {
	return &errorsPkg.PublicError{
		Title:      "Access token not found",
		Code:       "ACCESS_TOKEN_NOT_FOUND",
		HttpStatus: http.StatusUnauthorized,
	}
}

func accessTokenInvalidError() *errorsPkg.PublicError {
	return &errorsPkg.PublicError{
		Title:      "Access token is invalid",
		Code:       "ACCESS_TOKEN_INVALID",
		HttpStatus: http.StatusUnauthorized,
	}
}
