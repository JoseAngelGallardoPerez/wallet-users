package auth

import (
	"time"

	"github.com/Confialink/wallet-pkg-utils"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

const ReservePadding = "10s"

type AutologoutTTLResolver struct {
	sysSettings *syssettings.SysSettings
}

func NewAutologoutTTLResolver(sysSettings *syssettings.SysSettings) TokenTTLResolver {
	return &AutologoutTTLResolver{sysSettings}
}

// ResolveByTokenSubject calculates tokens lifetime
// if autologout is enabled a refresh token lifetime = autologout timeout + padding timeout + reserve padding time
// an access token lifetime = (autologout + padding timeout) / 2 + reserve padding time
// autologout timeout, padding time - default autologout settings
// reserve padding time - additional time to client be able to handle refresh or autologout
func (a *AutologoutTTLResolver) ResolveByTokenSubject(subject string) (time.Duration, error) {
	autologoutSettings, err := a.sysSettings.GetAutologoutSettings()
	defaultValues := map[string]string{
		ClaimAccessSub:  ClaimAccessTokenExp,
		ClaimRefreshSub: ClaimRefreshTokenExp,
	}

	if err != nil {
		return time.Duration(0), err
	}

	if !autologoutSettings.Enabled {
		if defaultValue, ok := defaultValues[subject]; ok {
			return utils.MustParseDuration(defaultValue), nil
		}

		return time.Duration(0), nil
	}

	if autologoutSettings.Timeout.Minutes() < 1.0 {
		autologoutSettings.Timeout = time.Duration(1 * time.Minute)
	}

	accessDuration := (autologoutSettings.Timeout+autologoutSettings.Padding)/2 + utils.MustParseDuration(ReservePadding)

	if subject == ClaimAccessSub {
		return accessDuration, nil
	}

	if subject == ClaimRefreshSub {
		return accessDuration + autologoutSettings.Timeout + autologoutSettings.Padding, nil
	}

	return time.Duration(0), nil
}

type FixedValueTTLResolver struct {
	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
}

func NewFixedValueTTLResolver(refreshTokenTTL string, accessTokenTTL string) *FixedValueTTLResolver {
	return &FixedValueTTLResolver{
		refreshTokenTTL: utils.MustParseDuration(refreshTokenTTL),
		accessTokenTTL:  utils.MustParseDuration(accessTokenTTL),
	}
}

func (f *FixedValueTTLResolver) ResolveByTokenSubject(subject string) (time.Duration, error) {
	switch subject {
	case ClaimRefreshSub:
		return f.refreshTokenTTL, nil
	case ClaimAccessSub:
		return f.accessTokenTTL, nil
	}

	return time.Duration(0), nil
}
