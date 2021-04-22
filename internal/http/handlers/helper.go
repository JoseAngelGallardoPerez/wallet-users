package handlers

import (
	"errors"
	"strconv"

	"github.com/Confialink/wallet-users/internal/db/models"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
)

// getCurrentUser retrieve current user from gin context
func GetCurrentUser(ctx *gin.Context) *userpb.User {
	user, ok := ctx.Get("_user")
	if !ok {
		return nil
	}
	return user.(*userpb.User)
}

// permission mw already get user by UID from DB. retrieve requested user from gin context
func GetRequestedUser(ctx *gin.Context) *models.User {
	user, ok := ctx.Get("_requested_user")
	if !ok {
		return nil
	}
	return user.(*models.User)
}

func getInt64Param(ctx *gin.Context, name string) (int64, error) {
	paramStr, isSet := ctx.Params.Get(name)
	if !isSet {
		return 0, errors.New("parameter is not set")
	}

	intParam, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid parameter passed expected integer value got \"" + paramStr + "\"")
	}

	return intParam, nil
}

func getUint64Param(ctx *gin.Context, name string) (uint64, error) {
	paramStr, isSet := ctx.Params.Get(name)
	if !isSet {
		return 0, errors.New("parameter is not set")
	}

	res, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid parameter passed expected integer value got \"" + paramStr + "\"")
	}

	return res, nil
}
