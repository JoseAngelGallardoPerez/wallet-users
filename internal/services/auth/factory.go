package auth

import (
	base "github.com/dgrijalva/jwt-go"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/config"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/jwt"
	"github.com/Confialink/wallet-users/internal/services/notifications"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

func TokenServiceFactory(
	configuration *config.Configuration,
	repository *repositories.TokenRepository,
	resolver TokenTTLResolver,
	logger log15.Logger,
) *TokenService {
	var (
		jwtService jwt.Service
		err        error
	)
	cfg := configuration.JWT

	switch signingMethod := cfg.SigningMethod.(type) {
	case *base.SigningMethodHMAC:
		jwtService = jwt.NewWithHMAC(cfg.Secret, signingMethod)
	case *base.SigningMethodECDSA:
		params := &jwt.Params{
			SigningMethod: signingMethod,
			SecretKeyPath: cfg.SecretKeyPath,
			PublicKeyPath: cfg.PublicKeyPath,
		}
		jwtService, err = jwt.NewWithECDSA(params)
	}
	if err != nil {
		logger.
			New("instance", "container", "method", "TemporaryTokensService").
			Crit("failed to create tokens service", "error", err)
		panic(err)
	}

	return NewTokenService(jwtService, repository, resolver)
}

func TemporaryTokensFactory(configuration *config.Configuration, logger log15.Logger) *TemporaryTokens {
	var (
		jwtService jwt.Service
		err        error
	)
	cfg := configuration.JWT

	switch signingMethod := cfg.SigningMethod.(type) {
	case *base.SigningMethodHMAC:
		jwtService = jwt.NewWithHMAC(cfg.Secret, signingMethod)
	case *base.SigningMethodECDSA:
		params := &jwt.Params{
			SigningMethod: signingMethod,
			SecretKeyPath: cfg.SecretKeyPath,
			PublicKeyPath: cfg.PublicKeyPath,
		}
		jwtService, err = jwt.NewWithECDSA(params)
	}
	if err != nil {
		logger.
			New("instance", "container", "method", "TemporaryTokensService").
			Crit("failed to create tokens service", "error", err)
		panic(err)
	}

	return NewTemporaryTokens(jwtService)
}

func BlockerFactory(
	blockedIpRepository *repositories.BlockedIpsRepository,
	failAuthAttemptRepository *repositories.FailAuthAttemptRepository,
	userRepository *repositories.UsersRepository,
	logger log15.Logger,
	notificationsService *notifications.Notifications,
	sysSettings *syssettings.SysSettings,
) *Blocker {
	return &Blocker{
		blockedIpRepository,
		failAuthAttemptRepository,
		userRepository,
		logger,
		&syssettings.LoginSecuritySettings{},
		notificationsService,
		sysSettings,
	}
}
