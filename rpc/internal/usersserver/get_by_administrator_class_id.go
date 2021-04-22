package usersserver

import (
	"context"
	"github.com/twitchtv/twirp"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"
)

// GetByUID returns user by uid
func (s *UsersHandlerServer) GetByAdministratorClassId(ctx context.Context, req *pb.Request) (res *pb.Response, err error) {
	users, err := s.Repository.GetUsersRepository().FindByClassId(req.ClassId)
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
