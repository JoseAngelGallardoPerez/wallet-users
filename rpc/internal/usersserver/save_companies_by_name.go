package usersserver

import (
	"context"
	"github.com/twitchtv/twirp"

	pb "github.com/Confialink/wallet-users/rpc/proto/users"
)

// SaveCompaniesByName save companies by name and returns
func (s *UsersHandlerServer) SaveCompaniesByName(ctx context.Context, req *pb.CompaniesNameRequest) (res *pb.CompaniesResponse, err error) {
	companies, err := s.Repository.GetCompanyRepository().SaveAndFindByNames(req.Names)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	responseCompanies := make([]*pb.Company, len(companies))
	for i, c := range companies {
		responseCompanies[i] = &pb.Company{
			ID:                c.ID,
			CompanyName:       c.CompanyName,
			CompanyType:       c.CompanyType,
			CompanyRole:       c.CompanyRole,
			DirectorFirstName: c.DirectorFirstName,
			DirectorLastName:  c.DirectorLastName,
		}
	}

	result := &pb.CompaniesResponse{
		Companies: responseCompanies,
	}

	return result, nil
}
