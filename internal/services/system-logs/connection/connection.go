package connection

import (
	"github.com/Confialink/wallet-users/internal/srvdiscovery"
	"net/http"

	pb "github.com/Confialink/wallet-logs/rpc/logs"
)

func GetSystemLogsClient() (pb.LogsService, error) {
	logsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameLogs)
	if err != nil {
		return nil, err
	}
	return pb.NewLogsServiceProtobufClient(logsUrl.String(), http.DefaultClient), nil
}
