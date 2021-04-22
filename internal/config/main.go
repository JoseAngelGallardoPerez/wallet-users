package config

import (
	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/inconshreveable/log15"
)

// Configuration is microservice config model
type Configuration struct {
	Server        *ServerConfiguration
	Database      *env_config.Db
	Cors          *env_config.Cors
	RPC           *RPCConfiguration
	JWT           *JwtConfiguration
	MessageBroker *MessageBroker
}

// Create a new config instance.
var conf *Configuration

// init initializes conf variables
func InitConfig(logger log15.Logger) (err error) {
	defaultConfigReader := env_config.NewReader("users")

	// Configure server.
	server := &ServerConfiguration{}
	if err = server.Init(); err != nil {
		return
	}

	// Configure database settings.
	database := defaultConfigReader.ReadDbConfig()

	// Configure cors settings.
	cors := defaultConfigReader.ReadCorsConfig()

	// Configure rpc settings.
	rpc := &RPCConfiguration{}
	if err = rpc.Init(); err != nil {
		return
	}

	jwt := &JwtConfiguration{}
	if err = jwt.Init(); err != nil {
		return
	}

	//conf.Server = cfg
	conf = &Configuration{
		Server:        server,
		Database:      database,
		Cors:          cors,
		RPC:           rpc,
		JWT:           jwt,
		MessageBroker: initMessageBrokerConfig(),
	}

	validateConfig(conf, logger)
	return
}

// GetConf returns the conf Configuration struct.
func GetConf() *Configuration {
	return conf
}

// GetServer returns server configuration
func (s *Configuration) GetServer() *ServerConfiguration {
	return s.Server
}

// GetDatabase returns database configuration
func (s *Configuration) GetDatabase() *env_config.Db {
	return s.Database
}

// GetCors returns cors configuration
func (s *Configuration) GetCors() *env_config.Cors {
	return s.Cors
}

// GetRPC returns RPC configuration
func (s *Configuration) GetRPC() *RPCConfiguration {
	return s.RPC
}

func validateConfig(cfg *Configuration, logger log15.Logger) {
	validator := env_config.NewValidator(logger)
	validator.ValidateCors(cfg.Cors, logger)
	validator.ValidateDb(cfg.Database, logger)
	validator.CriticalIfEmpty(cfg.Server.Port, "VELMIE_WALLET_USERS_SERVER_PORT", logger)
	validator.CriticalIfEmpty(cfg.RPC.UsersServerPort, "VELMIE_WALLET_USERS_RPC_PORT", logger)
}
