package workers

import (
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/inconshreveable/log15"
)

func newUpdateDormantUsers(
	repo *repositories.UsersRepository,
	tokenService *auth.TokenService,
	sysSettings *syssettings.SysSettings,
	logger log15.Logger,
) *updateDormantUsers {
	return &updateDormantUsers{
		repo,
		tokenService,
		sysSettings,
		logger.New("Worker", "updateDormantUsers"),
	}
}

func newRemoveInvalidTokens(repo *repositories.TokenRepository, service *auth.TokenService, logger log15.Logger) *removeInvalidTokens {
	return &removeInvalidTokens{
		tokenRepository: repo,
		tokenService:    service,
		logger:          logger.New("Worker", "removeInvalidTokens"),
	}
}
