package usersserver

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/Confialink/wallet-users/internal/services/users"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/inconshreveable/log15"
)

// UsersHandlerServer implements the Users service
type UsersHandlerServer struct {
	Repository        repositories.RepositoryInterface
	tokenService      *auth.TokenService
	tmpTokenService   *auth.TemporaryTokens
	sysSettings       *syssettings.SysSettings
	userService       *users.UserService
	userLoaderService *users.UserLoaderService
	logger            log15.Logger
}

func NewUserHandlerServer(
	repsitory repositories.RepositoryInterface,
	tokenService *auth.TokenService,
	tmpTokenService *auth.TemporaryTokens,
	sysSettings *syssettings.SysSettings,
	userService *users.UserService,
	userLoaderService *users.UserLoaderService,
	logger log15.Logger,
) *UsersHandlerServer {
	return &UsersHandlerServer{
		Repository:        repsitory,
		tokenService:      tokenService,
		tmpTokenService:   tmpTokenService,
		sysSettings:       sysSettings,
		userService:       userService,
		userLoaderService: userLoaderService,
		logger:            logger,
	}
}

func getResponseUser(user *models.User) *pb.User {
	groupId := uint64(0)
	if user.UserGroupId != nil {
		groupId = *user.UserGroupId
	}
	companyID := uint64(0)
	if user.CompanyID != nil {
		companyID = *user.CompanyID
	}
	var smsPhoneNumber string
	if user.SmsPhoneNumber != nil {
		smsPhoneNumber = *user.SmsPhoneNumber
	}
	result := &pb.User{
		UID:                    user.UID,
		Email:                  user.Email,
		Username:               user.Username,
		FirstName:              user.FirstName,
		LastName:               user.LastName,
		RoleName:               user.RoleName,
		GroupId:                groupId,
		CompanyName:            user.CompanyDetails.CompanyName,
		PhoneNumber:            user.PhoneNumber,
		SmsPhoneNumber:         smsPhoneNumber,
		ParentUID:              user.ParentId,
		CompanyID:              companyID,
		IsEmailConfirmed:       user.IsEmailConfirmed,
		IsPhoneNumberConfirmed: user.IsEmailConfirmed,
	}

	if classId, err := user.ClassId.Int64(); err == nil {
		result.AdministratorClassId = classId
	}
	return result
}
