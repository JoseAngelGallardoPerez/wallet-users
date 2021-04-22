package system_logs

import (
	"os"

	"github.com/Confialink/wallet-pkg-utils/recovery"
	"github.com/inconshreveable/log15"
)

func SystemLogsServiceFactory(logger log15.Logger) *SystemLogsService {
	return NewSystemLogsService(recovery.RecoveryWithWriter(os.Stdout), logger)
}
