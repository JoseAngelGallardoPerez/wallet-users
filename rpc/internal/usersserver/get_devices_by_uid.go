package usersserver

import (
	"context"
	"errors"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"
)

// GetDevicesByUID returns users by uid
func (s *UsersHandlerServer) GetDevicesByUID(ctx context.Context, req *pb.DevicesRequest) (res *pb.DevicesResponse, err error) {
	return nil, errors.New("this method id deprecated")
}
