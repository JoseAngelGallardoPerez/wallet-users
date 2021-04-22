package auth

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/auth"
)

const TmpAuthHeader = "X-Tmp-Auth"

type ResponseDto struct {
	Status  int
	Data    interface{}
	Headers map[string]string
}

type SignUpResponse struct {
	tokenService     *auth.TokenService
	tmpTokensService *auth.TemporaryTokens
}

func NewSignUpResponse(tokenService *auth.TokenService, tmpTokensService *auth.TemporaryTokens) *SignUpResponse {
	return &SignUpResponse{tokenService: tokenService, tmpTokensService: tmpTokensService}
}

// Make returns a response which depends on user's status.
// It returns an access and a refresh tokens for active users.
// And it returns a temporary token for other statuses
func (s *SignUpResponse) Make(user *models.User) (*ResponseDto, error) {
	if user.IsActive() {
		tokens, err := s.tokenService.IssueTokens(user, nil)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create tokens")
		}

		return &ResponseDto{
			Status:  http.StatusOK,
			Data:    tokens,
			Headers: nil,
		}, nil
	}

	tmpAuthToken, err := s.tmpTokensService.Issue(user)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create temporary auth token")
	}

	return &ResponseDto{
		Status:  http.StatusCreated,
		Data:    user,
		Headers: map[string]string{TmpAuthHeader: tmpAuthToken, "Access-Control-Expose-Headers": TmpAuthHeader},
	}, nil
}
