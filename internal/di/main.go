package di

import (
	"log"

	"github.com/inconshreveable/log15"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/dig"

	"github.com/Confialink/wallet-users/internal/config"
	"github.com/Confialink/wallet-users/internal/db"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	formconditions "github.com/Confialink/wallet-users/internal/http/form-conditions"
	"github.com/Confialink/wallet-users/internal/http/forms"
	"github.com/Confialink/wallet-users/internal/http/handlers"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/http/routes"
	httpAuth "github.com/Confialink/wallet-users/internal/http/services/auth"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/csv"
	"github.com/Confialink/wallet-users/internal/services/invites"
	messagebroker "github.com/Confialink/wallet-users/internal/services/message-broker"
	"github.com/Confialink/wallet-users/internal/services/users"
	"github.com/Confialink/wallet-users/internal/validators"
	"github.com/Confialink/wallet-users/rpc/cmd/server/usersserver"
)

var Container *dig.Container

func init() {
	Container = dig.New()

	providers := []interface{}{
		// *gorm.DB
		db.ConnectionFactory,
		// log15.Logger
		log15.New,
		// *gocron.Scheduler
		gocron.NewScheduler,
		routes.Api,
	}

	providers = append(providers, config.Providers()...)
	providers = append(providers, services.Providers()...)
	providers = append(providers, auth.Providers()...)
	providers = append(providers, csv.Providers()...)
	providers = append(providers, users.Providers()...)
	providers = append(providers, invites.Providers()...)
	providers = append(providers, repositories.Providers()...)
	providers = append(providers, usersserver.Providers()...)
	providers = append(providers, responses.Providers()...)
	providers = append(providers, handlers.Providers()...)
	providers = append(providers, messagebroker.Providers()...)
	providers = append(providers, forms.Providers()...)
	providers = append(providers, validators.Providers()...)
	providers = append(providers, formconditions.Providers()...)
	providers = append(providers, httpAuth.Providers()...)

	for _, provider := range providers {
		err := Container.Provide(provider)
		if err != nil {
			log.Fatal(err)
		}
	}
}
