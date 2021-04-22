package customization

import (
	"context"

	"github.com/pkg/errors"

	pb "github.com/Confialink/wallet-customization/rpc/proto"

	"github.com/Confialink/wallet-users/internal/services/customization/connection"
)

func GetCustomizationByKey(key string) (*pb.Customization, error) {
	client, err := connection.GetCustomizationClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GetByKey(context.Background(), &pb.Request{Key: key})
	if err != nil {
		return nil, err
	}

	if response.GetError() != nil {
		return nil, errors.New(response.GetError().Details)
	}

	return response.GetCustomization(), nil
}
