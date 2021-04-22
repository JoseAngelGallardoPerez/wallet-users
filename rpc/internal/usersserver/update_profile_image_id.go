package usersserver

import (
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"context"
)

func (s *UsersHandlerServer) UpdateProfileImageID(ctx context.Context, req *pb.UpdateProfileImageIDRequest) (res *pb.UpdateProfileImageIDResponse, err error) {
	res = &pb.UpdateProfileImageIDResponse{}
	repo := s.Repository.GetUsersRepository()
	user, err := repo.FindByUID(req.UID)
	if err != nil {
		return
	}

	user.ProfileImageID = &req.ImageID
	_, err = repo.Update(user)
	return
}
