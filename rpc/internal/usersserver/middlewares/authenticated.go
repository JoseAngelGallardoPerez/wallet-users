package middlewares

import (
	"fmt"
	"time"

	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
)

const defaultTrackDuration = "1h"

type Authenticated struct {
	repo   *repositories.UsersRepository
	logger log15.Logger
}

func (middleware *Authenticated) Call(user *models.User) {
	fmt.Println(isNeedUpdateLastActedAt(user.LastActedAt))
	if isNeedUpdateLastActedAt(user.LastActedAt) {
		middleware.updateLastActedAt(user)
	}
}

func (middleware *Authenticated) updateLastActedAt(user *models.User) {
	user.LastActedAt = time.Now()
	if _, err := middleware.repo.Update(user); err != nil {
		middleware.logger.Error("Can't update user", "error", err)
	}
}

func isNeedUpdateLastActedAt(currentValue time.Time) bool {
	lastValueWithDuration := currentValue.Add(getDuration())
	return lastValueWithDuration.Before(time.Now())
}

func getDuration() time.Duration {
	duration, _ := time.ParseDuration(defaultTrackDuration)
	return duration
}
