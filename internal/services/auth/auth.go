package auth

import (
	"fmt"

	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

type Auth struct {
	userRepo        *repositories.UsersRepository
	authBlocker     *Blocker
	sysSettings     *syssettings.SysSettings
	passwordService *services.Password
	tokenService    *TokenService
	logger          log15.Logger
}

func NewAuth(
	userRepo *repositories.UsersRepository,
	authBlocker *Blocker,
	sysSettings *syssettings.SysSettings,
	passwordService *services.Password,
	tokenService *TokenService,
	logger log15.Logger,
) *Auth {
	return &Auth{
		userRepo,
		authBlocker,
		sysSettings,
		passwordService,
		tokenService,
		logger,
	}
}

// it returns tokens by username and password
// we apply blocking user policy
func (s *Auth) LoginUser(user *models.User, userModel models.User, tokenOptions *TokenOptions, ip string) (*ExtendedTokensResponse, *responses.Error) {
	if err := s.passwordService.UserCheckPassword(userModel.Password, user.Password); err != nil {
		s.authBlocker.AddUserFailAttempt(userModel.Email, ip)
		return nil, responses.NewCommonErrorByCode(responses.CodeInvalidUsernamePassword, "Invalid username or password.")
	}

	tokens, errResp := s.issueTokens(user, tokenOptions)
	if errResp != nil {
		return nil, errResp
	}

	return tokens, nil
}

func (s *Auth) issueTokens(user *models.User, tokenOptions *TokenOptions) (*ExtendedTokensResponse, *responses.Error) {
	if errResp := s.checkMaintenanceMode(user); errResp != nil {
		return nil, errResp
	}

	if errResp := s.checkUserStatus(user); errResp != nil {
		return nil, errResp
	}

	tokens, err := s.tokenService.IssueTokens(user, tokenOptions)
	if err != nil {
		logger := s.logger.New("method", "issueTokens")
		logger.Error("failed to generate access token", "error", err)
		return nil, responses.NewCommonErrorByCode(responses.Unauthorized, "")
	}

	return &ExtendedTokensResponse{
		TokensResponse: *tokens,
		ChallengeName:  user.ChallengeName,
	}, nil
}

func (s *Auth) checkMaintenanceMode(user *models.User) *responses.Error {
	if !user.IsRoot() {
		maintenanceModeSettings, err := s.sysSettings.GetMaintenanceModeSettings()
		if err != nil {
			logger := s.logger.New("method", "checkMaintenanceMode")
			logger.Error("cannot retrieve MaintenanceModeSettings", "error", err)
			return responses.NewCommonErrorByCode(responses.MaintenanceMode, "")
		}

		if maintenanceModeSettings.Enabled {
			return responses.NewCommonErrorByCode(responses.MaintenanceMode, "")
		}
	}

	return nil
}

func (s *Auth) checkUserStatus(user *models.User) *responses.Error {
	if !user.IsActive() {
		switch user.Status {
		case models.StatusDormant:
			return responses.NewCommonErrorByCode(responses.CodeUserIsDormant, fmt.Sprintf("Sorry, user %s is dormant", user.Username))

		case models.StatusBlocked, models.StatusCanceled: // Returns a "423 StatusLocked" response
			return responses.NewCommonErrorByCode(responses.CodeUserIsNotActive, fmt.Sprintf("Sorry, user %s is inactive ", user.Username))

		case models.StatusPending:
			return responses.NewCommonErrorByCode(responses.CodeUserIsPending, fmt.Sprintf("Sorry, user %s is pending ", user.Username))

		default:
			return responses.NewCommonErrorByCode(responses.CodeUserIsNotActive, fmt.Sprintf("Sorry, user %s is inactive ", user.Username))

		}
	}

	return nil
}
