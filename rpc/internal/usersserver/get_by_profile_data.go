package usersserver

import (
	"context"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/twitchtv/twirp"
)

// GetByUsername returns users by username
func (s *UsersHandlerServer) GetByProfileData(ctx context.Context,
	req *pb.Request) (res *pb.Response, err error) {
	users, err := s.Repository.GetUsersRepository().FindByProfileData(req.Username, req.SearchColumns)
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
