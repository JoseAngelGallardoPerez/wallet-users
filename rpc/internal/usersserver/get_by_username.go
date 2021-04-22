package usersserver

import (
	"context"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/twitchtv/twirp"
)

// GetByUsername returns users by username
func (s *UsersHandlerServer) GetByUsername(ctx context.Context,
	req *pb.Request) (res *pb.Response, err error) {
	user, err := s.Repository.GetUsersRepository().FindByUsername(req.Username)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	result := &pb.Response{
		User: getResponseUser(user),
	}

	return result, nil
}
