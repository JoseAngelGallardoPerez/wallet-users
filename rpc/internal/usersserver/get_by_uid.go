package usersserver

import (
	"context"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/twitchtv/twirp"
)

// GetByUID returns user by uid
func (s *UsersHandlerServer) GetByUID(ctx context.Context, req *pb.Request) (res *pb.Response, err error) {
	user, err := s.Repository.GetUsersRepository().FindByUID(req.UID)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	result := &pb.Response{
		User: getResponseUser(user),
	}

	return result, nil
}
