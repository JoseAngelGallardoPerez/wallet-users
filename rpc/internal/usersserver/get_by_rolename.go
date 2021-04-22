package usersserver

import (
	"context"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/twitchtv/twirp"
)

// GetByRoleName returns users by rolename
func (s *UsersHandlerServer) GetByRoleName(ctx context.Context, req *pb.Request) (res *pb.Response, err error) {
	users, err := s.Repository.GetUsersRepository().FindByRoleName(req.RoleName)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	responseUsers := make([]*pb.User, len(users))
	for i, v := range users {
		responseUsers[i] = getResponseUser(v)
	}

	result := &pb.Response{
		Users: responseUsers,
	}

	return result, nil
}
