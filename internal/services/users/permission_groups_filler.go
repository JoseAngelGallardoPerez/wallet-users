package users

import (
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/permissions"
)

// PermissionGroupsFiller service for filling users by permission groups
type PermissionGroupsFiller struct {
	permissionsService *permissions.Permissions
}

// NewPermissionGroupsFiller returns new PermissionGroupsFiller
func NewPermissionGroupsFiller(permissionsService *permissions.Permissions) *PermissionGroupsFiller {
	return &PermissionGroupsFiller{permissionsService}
}

// FillUsers fills users by permission groups
func (f *PermissionGroupsFiller) FillUsers(users []*models.User) error {
	groupIds, err := getGroupIds(users)
	if err != nil {
		return err
	}
	groups := f.permissionsService.GetGroups(groupIds)
	fillUsersByGroups(users, groups)
	return nil
}

func fillUsersByGroups(users []*models.User, groups []*models.PermissionGroup) {
	for _, user := range users {
		user.PermissionGroup = getGroupByID(user.GetClassIDAsInt64(), groups)
	}
}

func getGroupByID(id int64, groups []*models.PermissionGroup) *models.PermissionGroup {
	if id == 0 {
		return nil
	}
	for _, group := range groups {
		if group.ID == id {
			return group
		}
	}
	return nil
}

func getGroupIds(users []*models.User) ([]int64, error) {
	ids := make(map[int64]interface{})
	for _, user := range users {
		id := user.GetClassIDAsInt64()
		if id == 0 {
			continue
		}
		ids[id] = nil
	}

	idsArray := make([]int64, 0)
	for k := range ids {
		idsArray = append(idsArray, k)
	}
	return idsArray, nil
}
