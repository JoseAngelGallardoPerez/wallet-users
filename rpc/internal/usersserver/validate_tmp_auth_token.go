package usersserver

import (
	"context"
	"github.com/dgrijalva/jwt-go"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"

	"github.com/Confialink/wallet-users/rpc/internal/usersserver/middlewares"
)

// ValidateTmpAuthToken validates temporary token and returns current user
func (s *UsersHandlerServer) ValidateTmpAuthToken(ctx context.Context, req *pb.Request) (res *pb.Response, err error) {
	token, err := s.tmpTokenService.Verify(req.TmpAuthToken)
	if err != nil {
		result := &pb.Response{
			Error: &pb.Error{
				Title:   "Token is not valid",
				Details: err.Error(),
			},
		}
		return result, err
	}

	claims := token.Claims.(jwt.MapClaims)
	user, err := s.Repository.GetUsersRepository().FindByUID(claims["uid"].(string))
	if err != nil {
		result := &pb.Response{
			Error: &pb.Error{
				Title:   "User not found",
				Details: err.Error(),
			},
		}
		return result, err
	}

	middlewares.NewAuthenticated(s.Repository.GetUsersRepository(), s.logger).Call(user)

	result := &pb.Response{
		User: getResponseUser(user),
	}

	return result, nil
}
