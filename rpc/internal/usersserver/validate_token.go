package usersserver

import (
	"github.com/Confialink/wallet-users/internal/services/auth"
	"context"
	"errors"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"

	"github.com/Confialink/wallet-users/rpc/internal/usersserver/middlewares"
)

// ValidateAccessToken validates token and returns current user
func (s *UsersHandlerServer) ValidateAccessToken(ctx context.Context, req *pb.Request) (res *pb.Response, err error) {
	_, err = s.tokenService.VerifyToken(req.AccessToken)
	if err != nil {

		result := &pb.Response{
			Error: &pb.Error{
				Title:   "Token is not valid",
				Details: err.Error(),
			},
		}
		return result, err
	}

	user, err := s.Repository.GetUsersRepository().FindUserByTokenAndSubject(req.AccessToken, auth.ClaimAccessSub)
	if err != nil {
		result := &pb.Response{
			Error: &pb.Error{
				Title:   "User not found",
				Details: err.Error(),
			},
		}
		return result, err
	}

	if user.RoleName != "root" {
		maintenanceModeSettings, err := s.sysSettings.GetMaintenanceModeSettings()
		if err != nil {
			result := &pb.Response{
				Error: &pb.Error{
					Title:   "Cannot get maintenance mode settings",
					Details: err.Error(),
				},
			}
			return result, err
		}

		if maintenanceModeSettings.Enabled {
			result := &pb.Response{
				Error: &pb.Error{
					Title: "Maintenance mode enabled",
				},
			}
			return result, errors.New("maintenance mode enabled")
		}
	}

	middlewares.NewAuthenticated(s.Repository.GetUsersRepository(), s.logger).Call(user)

	result := &pb.Response{
		User: getResponseUser(user),
	}

	return result, nil
}
