package connection

import (
	"github.com/Confialink/wallet-users/internal/srvdiscovery"
	"net/http"

	pb "github.com/Confialink/wallet-customization/rpc/proto"
)

func GetCustomizationClient() (pb.CustomizationHandler, error) {
	customizationUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameCustomization)
	if err != nil {
		return nil, err
	}
	return pb.NewCustomizationHandlerProtobufClient(customizationUrl.String(), http.DefaultClient), nil
}
