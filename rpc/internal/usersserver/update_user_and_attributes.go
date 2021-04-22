package usersserver

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
	"context"
)

// UpdateUserAndAttributes update user and attributes values
func (s *UsersHandlerServer) UpdateUserAndAttributes(ctx context.Context, req *pb.UserUpdateRequest,
) (res *pb.UserUpdateResponse, err error) {

	user := &models.User{
		UID:         req.UID,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
	}

	attributes := make(map[string]interface{})
	for _, attribute := range req.Attributes {
		attributes[attribute.Name] = attribute.Value
	}

	user.Attributes = attributes

	for _, address := range req.MailingAddresses {
		user.MailingAddresses = append(user.MailingAddresses, &models.Address{
			ID:                address.Id,
			UserID:            user.UID,
			Type:              models.AddressTypeMailing,
			CountryIsoTwo:     address.CountryIsoTwo,
			Region:            address.Region,
			City:              address.City,
			ZipCode:           address.ZipCode,
			Address:           address.Address,
			AddressSecondLine: address.AddressSecondLine,
			Name:              address.Name,
			PhoneNumber:       address.PhoneNumber,
			Description:       address.Description,
			Latitude:          &address.Latitude,
			Longitude:         &address.Longitude,
		})
	}

	for _, address := range req.PhysicalAddresses {
		user.PhysicalAddresses = append(user.PhysicalAddresses, &models.Address{
			ID:                address.Id,
			UserID:            user.UID,
			Type:              models.AddressTypePhysical,
			CountryIsoTwo:     address.CountryIsoTwo,
			Region:            address.Region,
			City:              address.City,
			ZipCode:           address.ZipCode,
			Address:           address.Address,
			AddressSecondLine: address.AddressSecondLine,
			Name:              address.Name,
			PhoneNumber:       address.PhoneNumber,
			Description:       address.Description,
			Latitude:          &address.Latitude,
			Longitude:         &address.Longitude,
		})
	}

	ts := s.Repository.GetUsersRepository().DB.Begin()
	err = s.userService.Update(user, ts)
	if err != nil {
		ts.Rollback()
		return nil, err
	}

	ts.Commit()

	return &pb.UserUpdateResponse{}, nil
}
