package services

import (
	"github.com/Confialink/wallet-users/internal/services/accounts"
	"github.com/Confialink/wallet-users/internal/services/files"
	"github.com/Confialink/wallet-users/internal/services/gdpr"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/permissions"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	system_logs "github.com/Confialink/wallet-users/internal/services/system-logs"
	"github.com/Confialink/wallet-users/internal/services/verification"
)

func Providers() []interface{} {
	return []interface{}{
		accounts.NewAccountsService,
		files.NewFilesService,
		permissions.NewPermissionsService,
		syssettings.NewSysSettings,
		syssettings.NewRpcClientFactory,
		system_logs.SystemLogsServiceFactory,
		NewPassword,
		notifications.NewRpcClientFactory,
		notifications.NewNotifications,
		verification.NewCreator,
		verification.NewValidator,
		gdpr.NewService,
	}
}
