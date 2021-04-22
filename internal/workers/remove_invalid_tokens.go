package workers

import (
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/repositories"
)

type removeInvalidTokens struct {
	tokenRepository *repositories.TokenRepository
	tokenService    *auth.TokenService
	logger          log15.Logger
}

func (w *removeInvalidTokens) execute() {
	refreshTokens, err := w.tokenRepository.FindTokensBySubject(auth.ClaimRefreshSub)
	if err != nil {
		w.logger.Error("can't find refresh tokens", err)
		return
	}

	for _, rt := range refreshTokens {
		_, err := w.tokenService.VerifyToken(rt.SignedString)
		if err != nil {
			accessToken, err := w.tokenRepository.FindAccessTokenByRefreshTokenId(rt.ID)
			if err != nil {
				_ = w.tokenRepository.Delete(rt)
			}

			_, err = w.tokenService.VerifyToken(accessToken.SignedString)
			if err != nil {
				_ = w.tokenRepository.Delete(accessToken)
				_ = w.tokenRepository.Delete(rt)
			}
		}
	}
}
