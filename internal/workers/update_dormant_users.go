package workers

import (
	"github.com/Confialink/wallet-users/internal/srvdiscovery"
	"context"
	"net/http"
	"time"

	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/inconshreveable/log15"

	notificationspb "github.com/Confialink/wallet-notifications/rpc/proto/notifications"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

type updateDormantUsers struct {
	repo         *repositories.UsersRepository
	tokenService *auth.TokenService
	sysSettings  *syssettings.SysSettings
	logger       log15.Logger
}

func (w *updateDormantUsers) execute() {
	duration, err := w.sysSettings.GetDormantDuration()
	if err != nil {
		w.logger.Error("Can't get dormant duration", "error", err)
		return
	}
	timeFilter := time.Now().Add(-duration).Format(time.RFC3339)
	params := list_params.NewListParamsFromQuery("", models.User{})
	params.AddFilter("lastActedAt", []string{timeFilter}, list_params.OperatorLt)
	params.AddFilter("roleName", []string{"client"})
	params.AddFilter("status", []string{models.StatusDormant}, list_params.OperatorNeq)
	params.Pagination.PageSize = 0
	users, err := w.repo.GetList(params)
	if err != nil {
		w.logger.Error("Can't get users", "error", err)
	}

	for _, v := range users {
		w.setDormantStatus(v)
	}
}

func (w *updateDormantUsers) setDormantStatus(user *models.User) {
	user.Status = models.StatusDormant
	if _, err := w.repo.Update(user); err != nil {
		w.logger.Error("Can't update user", "error", err)
	}

	if err := w.tokenService.RevokeUserTokens(user); err != nil {
		w.logger.Error("Can't log out user", "error", err)
	}

	client, err := w.getClient()
	if err != nil {
		w.logger.Error("failed to get pb client", "error", err)
	}

	_, err = client.Dispatch(context.Background(), &notificationspb.Request{
		EventName: "DormantProfileAdmin",
		TemplateData: &notificationspb.TemplateData{
			UserName: user.Username,
		},
	})

}

func (s *updateDormantUsers) getClient() (notificationspb.NotificationHandler, error) {
	notificationsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameNotifications)
	if nil != err {
		return nil, err
	}
	return notificationspb.NewNotificationHandlerProtobufClient(notificationsUrl.String(), http.DefaultClient), nil
}
