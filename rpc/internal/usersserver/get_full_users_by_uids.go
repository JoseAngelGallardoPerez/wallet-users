package usersserver

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/Confialink/wallet-pkg-utils/value"
	"github.com/twitchtv/twirp"

	"github.com/Confialink/wallet-users/internal/db/models"
	pb "github.com/Confialink/wallet-users/rpc/proto/users"
)

var allowedFields = []interface{}{
	"UID",
	"Email",
	"Username",
	"FirstName",
	"LastName",
	"PhoneNumber",
	"IsCorporate",
	"RoleName", "Status",
	"UserGroupId",
	"CompanyID",
	map[string][]interface{}{
		"UserGroup": {"Name"},
	},
	map[string][]interface{}{
		"CompanyDetails": {"CompanyName"},
	},
}

// GetFullUsersByUIDs returns users by uid
func (s *UsersHandlerServer) GetFullUsersByUIDs(ctx context.Context, req *pb.RequestFullUsersByUIDs,
) (res *pb.FullUsersResponse, err error) {
	repo := s.Repository.GetUsersRepository()

	params := list_params.NewListParams()
	params.ObjectType = reflect.TypeOf(models.User{})

	params.AllowSelectFields(allowedFields)

	params.AddFilter("uid", req.UIDs, list_params.OperatorIn)
	if len(req.Fields) > 0 {
		transformedFields := transformSelectedFields(req.Fields)
		if len(transformedFields.NestedFields) > 0 {
			for _, v := range transformedFields.NestedFields {
				params.Includes.AddIncludes(v.PropName)
			}
		}
		params.SelectFields(transformedFields.ToInterfaceArray())
	}
	users, err := repo.GetList(params)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &pb.FullUsersResponse{
		FullUsers: usersToProtoFullUsers(users),
	}, nil
}

func usersToProtoFullUsers(users []*models.User) []*pb.FullUser {
	result := make([]*pb.FullUser, len(users))
	for i, v := range users {
		result[i] = userToProtoFullUser(v)
	}
	return result
}

func userToProtoFullUser(user *models.User) *pb.FullUser {
	return &pb.FullUser{
		Uid:         user.UID,
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsCorporate: value.FromBool(user.IsCorporate),
		RoleName:    user.RoleName,
		Status:      user.Status,
		UserGroupId: value.FromUint64(user.UserGroupId),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339Nano),
		UserDetails: userDetailsToFullUserDetails(&user.UserDetails),
		//PhysicalAdress:  userPhysicalAdressToFullPhysicalAdress(&user.PhysicalAdress),
		//BenificialOwner: userBenificialOwnerToFullBenificialOwner(&user.BenificialOwner),
		UserGroup:      userGroupToFullUserGroup(user.UserGroup),
		CompanyDetails: companyToFullCompany(&user.CompanyDetails),
	}
}

func userGroupToFullUserGroup(group *models.UserGroup) *pb.UserGroup {
	if group == nil {
		return nil
	}

	return &pb.UserGroup{
		Id:          group.ID,
		Name:        group.Name,
		Description: group.Description,
	}
}

func userDetailsToFullUserDetails(details *models.UserDetails) *pb.UserDetails {
	if details == nil {
		return nil
	}

	return &pb.UserDetails{
		ClassId:                  string(details.ClassId),
		CountryOfResidenceIso2:   details.CountryOfResidenceIsoTwo,
		CountryOfCitizenshipIso2: details.CountryOfCitizenshipIsoTwo,
		DocumentType:             details.GetDocumentType(),
		DocumentPersonalId:       details.DocumentPersonalId,
		Fax:                      details.Fax,
		HomePhoneNumber:          details.HomePhoneNumber,
		InternalNotes:            details.InternalNotes,
		OfficePhoneNumber:        details.OfficePhoneNumber,
		Position:                 details.Position,
	}
}
func companyToFullCompany(company *models.Company) *pb.Company {
	if company == nil {
		return nil
	}

	return &pb.Company{
		ID:                company.ID,
		CompanyName:       company.CompanyName,
		CompanyType:       company.CompanyType,
		CompanyRole:       company.CompanyRole,
		DirectorFirstName: company.DirectorFirstName,
		DirectorLastName:  company.DirectorLastName,
	}
}

func transformSelectedFields(fields []string) *list_params.Fields {
	innerStructs := []string{"UserDetails", "PhysicalAdress", "MailingAddress", "BenificialOwner"}
	for i, field := range fields {
		for _, innerStruct := range innerStructs {
			if strings.HasPrefix(field, innerStruct) {
				fields[i] = strings.TrimPrefix(field, innerStruct+".")
			}
		}
	}

	fieldsObject := list_params.StringsArrayToFields(fields)
	return &fieldsObject
}
