package workers

import (
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/inconshreveable/log15"
	"github.com/jasonlvhit/gocron"
)

func Start(
	scheduler *gocron.Scheduler,
	repo *repositories.UsersRepository,
	tokenService *auth.TokenService,
	sysSettings *syssettings.SysSettings,
	logger log15.Logger,
) {
	register(scheduler, repo, tokenService, sysSettings, logger)
	scheduler.Start()
	log15.New().Info("Scheduler is started")
}

func register(
	scheduler *gocron.Scheduler,
	repo *repositories.UsersRepository,
	tokenService *auth.TokenService,
	sysSettings *syssettings.SysSettings,
	logger log15.Logger,
) {
	scheduler.Every(1).Hour().Do(newUpdateDormantUsers(repo, tokenService, sysSettings, logger).execute)
	// scheduler.Every(24).Hour().Do(newRemoveInvalidTokens().execute)
}
