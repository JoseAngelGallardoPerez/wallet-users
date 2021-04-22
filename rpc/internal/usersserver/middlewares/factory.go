package middlewares

import (
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/inconshreveable/log15"
)

func NewAuthenticated(repo *repositories.UsersRepository, logger log15.Logger) *Authenticated {
	return &Authenticated{
		repo,
		logger.New("RPCMiddleware", "Authenticated"),
	}
}
