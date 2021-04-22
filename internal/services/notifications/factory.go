package notifications

import (
	"net/http"

	pb "github.com/Confialink/wallet-notifications/rpc/proto/notifications"

	"github.com/Confialink/wallet-users/internal/srvdiscovery"
)

type ClientFactory interface {
	// Creates a new HTTP client
	NewClient() (pb.NotificationHandler, error)
}

type RpcClientFactory struct {
}

func NewRpcClientFactory() ClientFactory {
	return &RpcClientFactory{}
}

func (s *RpcClientFactory) NewClient() (pb.NotificationHandler, error) {
	notificationsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameNotifications)
	if nil != err {
		return nil, err
	}
	return pb.NewNotificationHandlerProtobufClient(notificationsUrl.String(), http.DefaultClient), nil
}
