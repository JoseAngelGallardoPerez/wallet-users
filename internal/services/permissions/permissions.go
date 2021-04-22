package permissions

import (
	"github.com/Confialink/wallet-users/internal/srvdiscovery"
	"context"
	"log"
	"net/http"

	pbPermissions "github.com/Confialink/wallet-permissions/rpc/permissions"
	"github.com/Confialink/wallet-pkg-acl"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
)

const (
	ViewUserProfilesKey                 = "view_user_profiles"
	CreateUserProfilesKey               = "create_profile"
	ModifyUserProfilesKey               = "modify_user_profiles"
	ViewAdminProfilesKey                = "view_admin_profiles"
	CreateAdminProfilesKey              = "create_admin_profiles"
	ModifyAdminProfilesKey              = "modify_admin_profiles"
	CreateAccounts                      = "create_accounts"
	CreateAccountsWithInitialBalance    = "create_accounts_with_initial_balance"
	CreateCards                         = "create_cards"
	ModifyCards                         = "modify_cards"
	InitiateExecuteUserTransfers        = "initiate_execute_user_transfers"
	ViewUserReports                     = "view_user_reports"

	ViewSettings   = "view_settings"
	ModifySettings = "modify_settings"
	CreateSettings = "create_settings"
	RemoveSettings = "remove_settings"
)

type Permissions struct {
	checker         pbPermissions.PermissionChecker
	usersRepository *repositories.UsersRepository
	logger          log15.Logger
}

// NewPermissionsService creates a new Permissions instance.
func NewPermissionsService(usersRepository *repositories.UsersRepository, logger log15.Logger) *Permissions {
	return &Permissions{usersRepository: usersRepository, logger: logger}
}

// CanCreateProfile checks if can create new profile
func (p *Permissions) CanCreateProfile(uid string, requestedUser *models.User) bool {
	if requestedUser.IsAdmin() || requestedUser.IsRoot() {
		return p.CanCreateAdminProfile(uid, requestedUser)
	}
	return p.CanCreateUserProfile(uid)
}

// CanCreateAdminProfile checks if can create new admin profile
func (p *Permissions) CanCreateAdminProfile(uid string, requestedUser *models.User) bool {
	currentUser, err := p.usersRepository.FindByUID(uid)
	if err != nil {
		logger := p.logger.New("method", "CanCreateAdminProfile")
		logger.Error("failed to retrieve user by uid", "error", err, "uid", uid)
		return false
	}

	if acl.RolesHelper.FromName(currentUser.RoleName) < acl.RolesHelper.FromName(requestedUser.RoleName) {
		return false
	}

	return p.CheckPermission(uid, CreateAdminProfilesKey)
}

// CanCreateUserProfile checks if can create new user profile
func (p *Permissions) CanCreateUserProfile(uid string) bool {
	return p.CheckPermission(uid, CreateUserProfilesKey)
}

// CanViewProfile checks if user can view exists profile
func (p *Permissions) CanViewProfile(uid string, requestedUser *models.User) bool {
	if requestedUser.IsAdmin() || requestedUser.IsRoot() {
		return p.CanViewAdminProfile(uid)
	}
	return p.CanViewUserProfile(uid)
}

// CanViewAdminProfile checks if can view admin profile
func (p *Permissions) CanViewAdminProfile(uid string) bool {
	return p.CheckPermission(uid, ViewAdminProfilesKey)
}

// CanViewUserProfile checks if can view user profile
func (p *Permissions) CanViewUserProfile(uid string) bool {
	return p.CheckPermission(uid, ViewUserProfilesKey)
}

func (p *Permissions) CanViewShortUserProfiles(uid string) bool {
	permissions := []string{
		CreateAccounts,
		CreateAccountsWithInitialBalance,
		CreateCards,
		ModifyCards,
		InitiateExecuteUserTransfers,
		ViewUserReports,
		ViewAdminProfilesKey,
		ViewUserProfilesKey,
	}
	return p.CheckOneOfPermissions(uid, permissions)
}

// CanUpdateProfile checks if can update exists profile
func (p *Permissions) CanUpdateProfile(uid string, requestedUser *models.User) bool {
	if requestedUser.UID == uid {
		return true
	}

	if requestedUser.IsAdmin() || requestedUser.IsRoot() {
		return p.CanUpdateAdminProfile(uid, requestedUser)
	}
	return p.CanUpdateUserProfile(uid)
}

// CanUpdateAdminProfile checks if can update admin profile
func (p *Permissions) CanUpdateAdminProfile(uid string, requestedUser *models.User) bool {
	currentUser, err := p.usersRepository.FindByUID(uid)
	if err != nil {
		logger := p.logger.New("method", "CanUpdateAdminProfile")
		logger.Error("failed to retrieve user by uid", "error", err, "uid", uid)
		return false
	}

	if acl.RolesHelper.FromName(currentUser.RoleName) < acl.RolesHelper.FromName(requestedUser.RoleName) {
		return false
	}

	return p.CheckPermission(uid, ModifyAdminProfilesKey)
}

// CanUpdateUserProfile checks if can update user profile
func (p *Permissions) CanUpdateUserProfile(uid string) bool {
	return p.CheckPermission(uid, ModifyUserProfilesKey)
}

func (p *Permissions) CanViewSettings(uid string) bool {
	return p.CheckPermission(uid, ViewSettings)
}

func (p *Permissions) CanModifySettings(uid string) bool {
	return p.CheckPermission(uid, ModifySettings)
}

func (p *Permissions) CanCreateSettings(uid string) bool {
	return p.CheckPermission(uid, CreateSettings)
}

func (p *Permissions) CanRemoveSettings(uid string) bool {
	return p.CheckPermission(uid, RemoveSettings)
}

// CheckPermission checks permission
func (p *Permissions) CheckPermission(uid, actionKey string) bool {
	logger := p.logger.New("method", "CheckPermission")

	user, err := p.usersRepository.FindByUID(uid)
	if err != nil {
		logger.Error("failed to retrieve user by uid", "error", err, "uid", uid)
		return false
	}

	if user.RoleName == models.RoleRoot {
		return true
	}

	request := &pbPermissions.PermissionReq{
		UserId:    uid,
		ActionKey: actionKey,
	}

	checker, err := p.getChecker()
	if err != nil {
		return false
	}

	response, err := checker.Check(context.Background(), request)
	if nil != err {
		log.Printf("Failed to get permission response: %v", err)
		return false
	}

	if !response.IsAllowed {
		return false
	}
	return true
}

// return true if one of permissions is granted
func (p *Permissions) CheckOneOfPermissions(uid string, actionKeys []string) bool {
	logger := p.logger.New("method", "CheckPermission")

	user, err := p.usersRepository.FindByUID(uid)
	if err != nil {
		logger.Error("failed to retrieve user by uid", "error", err, "uid", uid)
		return false
	}

	if user.RoleName == models.RoleRoot {
		return true
	}

	request := &pbPermissions.PermissionsReq{
		UserId:     uid,
		ActionKeys: actionKeys,
	}

	checker, err := p.getChecker()
	if err != nil {
		return false
	}

	response, err := checker.CheckAll(context.Background(), request)
	if nil != err {
		logger.Error("failed to get permission response", "error", err, "uid", uid)
		return false
	}

	for _, perm := range response.Permissions {
		if perm.IsAllowed {
			return true
		}
	}

	return false
}

// GetGroups returns list of groups from permissions service
func (p *Permissions) GetGroups(ids []int64) []*models.PermissionGroup {
	checker, err := p.getChecker()
	if err != nil {
		p.logger.Error("Unable to get permissions checker", "error", err)
		return nil
	}
	resp, err := checker.GetGroupsByIds(context.Background(), &pbPermissions.GroupIdsReq{Ids: ids})
	if err != nil {
		p.logger.Error("Unable to get permission groups", "error", err)
		return nil
	}
	groups := make([]*models.PermissionGroup, len(resp.Groups))
	for i, v := range resp.Groups {
		groups[i] = p.transformGroup(v)
	}

	return groups
}

func (p *Permissions) transformGroup(group *pbPermissions.Group) *models.PermissionGroup {
	return &models.PermissionGroup{
		ID:          group.Id,
		Name:        group.Name,
		Description: group.Description,
	}
}

// getChecker return permissions checker
func (p *Permissions) getChecker() (pbPermissions.PermissionChecker, error) {
	if p.checker == nil {
		permissionsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNamePermissions)
		if err != nil {
			log.Printf("Failed to get permissions rpc url: %v", err)
			return nil, err
		}
		p.checker = pbPermissions.NewPermissionCheckerProtobufClient(permissionsUrl.String(), http.DefaultClient)
	}
	return p.checker, nil
}
