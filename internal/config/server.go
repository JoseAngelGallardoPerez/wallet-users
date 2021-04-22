package config

import (
	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/Confialink/wallet-pkg-env_mods"
)

// ServerConfiguration is server config model
type ServerConfiguration struct {
	Port string
	Env  string
}

// GetPort returns server port
func (s *ServerConfiguration) GetPort() string {
	return s.Port
}

// GetEnv returns env mode
func (s *ServerConfiguration) GetEnv() string {
	return s.Env
}

// Init initializes enviroment variables
func (s *ServerConfiguration) Init() error {
	s.Port = env_config.Env("VELMIE_WALLET_USERS_SERVER_PORT", "")
	s.Env = env_config.Env("ENV", env_mods.Development)
	return nil
}
