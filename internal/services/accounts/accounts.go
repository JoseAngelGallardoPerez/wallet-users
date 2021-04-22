package accounts

import (
	"context"
	"net/http"

	"github.com/Confialink/wallet-users/internal/srvdiscovery"

	"github.com/inconshreveable/log15"

	pb "github.com/Confialink/wallet-accounts/rpc/accounts"

	"github.com/Confialink/wallet-users/internal/db/models"
)

type AccountsService struct {
	logger log15.Logger
}

func NewAccountsService(logger log15.Logger) *AccountsService {
	return &AccountsService{logger}
}

func (s *AccountsService) UserHasCardsOrAccounts(uid string) (resp *pb.UserHasCardsOrAccountsResp, err error) {
	processor, err := s.processor()
	if err != nil {
		return nil, err
	}

	request := pb.UserHasCardsOrAccountsReq{Uid: uid}
	resp, err = processor.UserHasCardsOrAccountsBy(context.Background(), &request)
	if err != nil {
		s.logger.Error("cannot find accounts or cards", "err", err)
		return resp, err
	}

	return resp, nil
}

func (s *AccountsService) GenerateAccount(user *models.User) error {
	processor, err := s.processor()
	if err != nil {
		return err
	}

	request := pb.GenerateAccountReq{Uid: user.UID, CurrencyCode: "EUR"}
	_, err = processor.GenerateAccount(context.Background(), &request)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountsService) processor() (pb.AccountsProcessor, error) {
	accountsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameAccounts)
	if err != nil {
		s.logger.Error("failed to connect to accounts", "err", err)
		return nil, err
	}

	return pb.NewAccountsProcessorProtobufClient(accountsUrl.String(), http.DefaultClient), nil
}
