package usersserver

import (
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"context"
)

// GetFullUsersByUIDs returns user and attributes by uid
func (s *UsersHandlerServer) GetUserAndAttributes(ctx context.Context, req *pb.UserGetRequest,
) (res *pb.UserGetResponse, err error) {

	user, err := s.Repository.GetUsersRepository().FindByUID(req.UID)
	if err != nil {
		return nil, err
	}

	err = s.userLoaderService.LoadUserCompletely(user)
	if err != nil {
		return nil, err
	}

	res = &pb.UserGetResponse{
		UID:         user.UID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}

	for index, value := range user.Attributes {
		res.Attributes = append(res.Attributes, &pb.Attribute{
			Name:  index,
			Value: value.(string),
		})
	}

	for _, address := range user.MailingAddresses {
		res.MailingAddresses = append(res.MailingAddresses, &pb.Address{
			Id:                address.ID,
			CountryIsoTwo:     address.CountryIsoTwo,
			Region:            address.Region,
			City:              address.City,
			ZipCode:           address.ZipCode,
			Address:           address.Address,
			AddressSecondLine: address.AddressSecondLine,
			Name:              address.Name,
			PhoneNumber:       address.PhoneNumber,
			Description:       address.Description,
			Latitude:          *address.Latitude,
			Longitude:         *address.Longitude,
		})
	}

	for _, address := range user.PhysicalAddresses {
		res.PhysicalAddresses = append(res.PhysicalAddresses, &pb.Address{
			Id:                address.ID,
			CountryIsoTwo:     address.CountryIsoTwo,
			Region:            address.Region,
			City:              address.City,
			ZipCode:           address.ZipCode,
			Address:           address.Address,
			AddressSecondLine: address.AddressSecondLine,
			Name:              address.Name,
			PhoneNumber:       address.PhoneNumber,
			Description:       address.Description,
			Latitude:          *address.Latitude,
			Longitude:         *address.Longitude,
		})
	}

	return res, nil
}
