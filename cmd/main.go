package main

import (
	"log"

	"github.com/Confialink/wallet-users/internal/commands"
	auth2 "github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/validators"
	"github.com/Confialink/wallet-users/internal/workers"
	"github.com/jasonlvhit/gocron"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/syssettings"

	"github.com/Confialink/wallet-users/rpc/cmd/server/usersserver"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-pkg-env_mods"
	"github.com/Confialink/wallet-users/internal/config"
	"github.com/Confialink/wallet-users/internal/di"
	"github.com/Confialink/wallet-users/internal/http/forms"
)

// main: main function
func main() {
	c := di.Container

	var cfg *config.Configuration
	var pbServer *usersserver.UsersServer
	var logger log15.Logger
	err := c.Invoke(func(
		config *config.Configuration,
		pb *usersserver.UsersServer,
		loggerDep log15.Logger,

		scheduler *gocron.Scheduler,
		usersRepo *repositories.UsersRepository,
		securityQuestionsRepo *repositories.SecurityQuestionRepository,
		userGroupsRepo *repositories.UserGroupsRepository,
		sysSettings *syssettings.SysSettings,
		passwordService *services.Password,
		tokenService *auth2.TokenService,
		formBuilder *forms.Factory,
		engineValidator *validator.Validate,
	) {
		cfg = config
		pbServer = pb
		gin.SetMode(env_mods.GetMode(config.Server.Env))
		logger = loggerDep

		validators.Register(usersRepo, securityQuestionsRepo, userGroupsRepo, sysSettings, passwordService, engineValidator)

		createRootUserCommand := commands.Init(c)
		commands.AddCommand(createRootUserCommand)
		commands.Run()

		workers.Start(scheduler, usersRepo, tokenService, sysSettings, logger)
		if err := formBuilder.InitForms(); err != nil {
			log.Fatal("cannot initialize forms: " + err.Error())
		}
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Starting API on port: %s", cfg.Server.Port)

	// Start proto buf server
	go pbServer.Init()

	// Start gin server
	err = c.Invoke(func(apiRoutes *gin.Engine) {
		err = apiRoutes.Run(":" + cfg.Server.Port)
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}

}
