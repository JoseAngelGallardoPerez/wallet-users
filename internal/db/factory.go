package db

import (
	"github.com/Confialink/wallet-pkg-env_mods"
	"github.com/Confialink/wallet-users/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"
)

func ConnectionFactory(configuration *config.Configuration, logger log15.Logger) *gorm.DB {
	c, err := CreateConnection(configuration.GetDatabase())
	logger = logger.New("ConnectionFactory")
	// defer database.Close()
	if nil != err {
		logger.Crit("Could not connect to DB", err)
	}

	if env_mods.GetMode(configuration.GetServer().GetEnv()) == gin.DebugMode {
		c.LogMode(true)
	}
	return c
}
