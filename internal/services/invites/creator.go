package invites

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/inconshreveable/log15"
	"github.com/jinzhu/gorm"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
)

const (
	minCodeLength = 10
	days          = 7
)

type Creator struct {
	db                *gorm.DB
	logger            log15.Logger
	invitesRepository *repositories.InvitesRepository
}

func NewCreator(
	db *gorm.DB,
	logger log15.Logger,
	invitesRepository *repositories.InvitesRepository,
) *Creator {
	return &Creator{
		db,
		logger,
		invitesRepository,
	}
}

func (c *Creator) Call(to, userId string, tx *gorm.DB) (*models.Invite, errors.TypedError) {
	if tx != nil {
		c.invitesRepository = c.invitesRepository.WrapContext(tx)
	}
	model, err := c.create(to, userId)
	if err != nil {
		return nil, &errors.PrivateError{OriginalError: err}
	}
	return model, nil
}

// Create creates a new invite.
func (c *Creator) create(to, userUID string) (*models.Invite, error) {
	code, err := c.generateRandomStringURLSafe(minCodeLength)
	if err != nil {
		panic(err)
	}

	invite := &models.Invite{
		Code:    code,
		Uses:    0,
		UserUID: userUID,
	}

	invite.RestrictUsageTo(to).CanBeUsedOnce().ExpiresIn(days)

	if err := c.invitesRepository.Create(invite); err != nil {
		return nil, &errors.PrivateError{OriginalError: err}
	}
	return invite, nil
}

// generateRandomStringURLSafe returns a URL-safe, base64 encoded securely generated random string.
func (c *Creator) generateRandomStringURLSafe(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}
