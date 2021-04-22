package system_logs

import (
	"context"

	pb "github.com/Confialink/wallet-logs/rpc/logs"
	"github.com/Confialink/wallet-users/internal/services/system-logs/connection"
)

type logsServiceWrap struct {
	logsService pb.LogsService
}

func newLogsServiceWrap() *logsServiceWrap {
	return &logsServiceWrap{}
}

func (self *logsServiceWrap) createLog(
	subject string,
	userId string,
	logTime string,
	dataTitle string,
	data []byte,
) {
	resp, err := self.systemLogger().CreateLog(context.Background(), &pb.CreateLogReq{
		Subject:    subject,
		UserId:     userId,
		LogTime:    logTime,
		DataTitle:  dataTitle,
		DataFields: data,
	})
	if err != nil {
		return
	}
	if resp.Error != nil {
		return
	}
}

func (self *logsServiceWrap) systemLogger() pb.LogsService {
	if self.logsService == nil {
		var err error
		self.logsService, err = connection.GetSystemLogsClient()
		if err != nil {
			return nil
		}
	}
	return self.logsService
}
