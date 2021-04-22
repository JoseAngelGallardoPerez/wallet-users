package config

import (
	"github.com/Confialink/wallet-pkg-env_config"
)

// RPCConfiguration is rpc config model
type RPCConfiguration struct {
	UsersServerPort string
}

// GetUsersServerPort returns rpc port for userserver
func (s *RPCConfiguration) GetUsersServerPort() string {
	return s.UsersServerPort
}

// Init initializes enviroment variables
func (s *RPCConfiguration) Init() error {
	s.UsersServerPort = env_config.Env("VELMIE_WALLET_USERS_RPC_PORT", "")
	return nil
}
