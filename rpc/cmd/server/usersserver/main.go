package usersserver

import (
	"github.com/Confialink/wallet-users/internal/services/users"
	"fmt"
	"github.com/inconshreveable/log15"
	"net/http"

	"github.com/Confialink/wallet-users/internal/config"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	server "github.com/Confialink/wallet-users/rpc/internal/usersserver"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
)

// UsersServer implements the Users service
type UsersServer struct {
	Repository        repositories.RepositoryInterface
	tokenService      *auth.TokenService
	tmpTokenService   *auth.TemporaryTokens
	sysSettings       *syssettings.SysSettings
	userService       *users.UserService
	userLoaderService *users.UserLoaderService
	logger            log15.Logger
}

func NewUsersServer(
	repsitory repositories.RepositoryInterface,
	tokenService *auth.TokenService,
	tmpTokenService *auth.TemporaryTokens,
	sysSettings *syssettings.SysSettings,
	userService *users.UserService,
	userLoaderService *users.UserLoaderService,
	logger log15.Logger,
) *UsersServer {
	return &UsersServer{
		Repository:        repsitory,
		tokenService:      tokenService,
		tmpTokenService:   tmpTokenService,
		sysSettings:       sysSettings,
		userService:       userService,
		userLoaderService: userLoaderService,
		logger:            logger,
	}
}

// Init initializes users rpc server
func (s *UsersServer) Init() {
	// Retrieve config options.
	conf := config.GetConf()

	hs := server.NewUserHandlerServer(s.Repository, s.tokenService, s.tmpTokenService, s.sysSettings, s.userService, s.userLoaderService, s.logger)

	twirpHandler := pb.NewUserHandlerServer(hs, nil)

	mux := http.NewServeMux()
	mux.Handle(pb.UserHandlerPathPrefix, twirpHandler)

	go http.ListenAndServe(fmt.Sprintf(":%s", conf.RPC.GetUsersServerPort()), mux)
}
