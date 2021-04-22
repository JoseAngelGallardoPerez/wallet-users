package syssettings

import (
	"net/http"

	pb "github.com/Confialink/wallet-settings/rpc/proto/settings"

	"github.com/Confialink/wallet-users/internal/srvdiscovery"
)

type ClientFactory interface {
	// Creates a new HTTP client
	NewClient() (pb.SettingsHandler, error)
}

type RpcClientFactory struct {
}

func NewRpcClientFactory() ClientFactory {
	return &RpcClientFactory{}
}

// NewClient creates a new HTTP client
func (s *RpcClientFactory) NewClient() (pb.SettingsHandler, error) {
	settingsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameSettings)
	if err != nil {
		return nil, err
	}
	return pb.NewSettingsHandlerProtobufClient(settingsUrl.String(), http.DefaultClient), nil
}
