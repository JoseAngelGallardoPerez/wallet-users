package handlers

import (
	"fmt"
	"strings"

	"github.com/Confialink/wallet-pkg-list_params"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/services/users"
)

// HandlerParams contains methods with params for handlers
type HandlerParams struct {
	permissionGroupsFiller *users.PermissionGroupsFiller
}

// NewHandlerParams returns new params for handler
func NewHandlerParams(permissionGroupsFiller *users.PermissionGroupsFiller) *HandlerParams {
	return &HandlerParams{permissionGroupsFiller}
}

func (p *HandlerParams) userProfilesCsv(query string) *list_params.ListParams {
	params := list_params.NewListParamsFromQuery(query, models.User{})
	params.Pagination.PageSize = 0
	userProfilesCsvIncludes(params)
	profilesCsvFilters(params)
	profilesCsvSortings(params)
	return params
}
func (p *HandlerParams) adminProfilesCsv(query string) *list_params.ListParams {
	params := list_params.NewListParamsFromQuery(query, models.User{})
	params.Pagination.PageSize = 0
	profilesCsvFilters(params)
	profilesCsvSortings(params)
	p.adminProfilesCsvIncludes(params)
	return params
}

func (p *HandlerParams) adminProfilesCsvIncludes(params *list_params.ListParams) {
	params.AllowIncludes([]string{"permissionGroup"})
	params.AddCustomIncludes("permissionGroup", func(records []interface{}) error {
		users := make([]*models.User, len(records))
		for i, v := range records {
			users[i] = v.(*models.User)
		}
		return p.permissionGroupsFiller.FillUsers(users)
	})
	params.Includes.AddIncludes("permissionGroup")
}

func userProfilesCsvIncludes(params *list_params.ListParams) {
	params.AllowIncludes([]string{"userGroup"})
	params.Includes.AddIncludes("userGroup")
}

func profilesCsvFilters(params *list_params.ListParams) {
	params.AllowFilters([]string{
		"query",
		"status",
		"user_group_id",
		"date_from",
		"date_to",
		"role_name",
		list_params.FilterIn("role_name"),
	})
	params.AddCustomFilter("query", queryFilter)
	params.AddCustomFilter("date_from", list_params.DateFromFilter("created_at"))
	params.AddCustomFilter("date_to", list_params.DateToFilter("created_at"))
}

func queryFilter(inputValues []string, params *list_params.ListParams) (
	dbConditionPart string, dbValues interface{},
) {
	columns := []string{"uid", "email", "username", "first_name", "last_name"}
	value := inputValues[0]

	conditions := make([]string, len(columns))
	values := make([]string, len(columns))
	for i, column := range columns {
		conditions[i] = fmt.Sprintf("`%s` LIKE ?", column)
		values[i] = fmt.Sprintf("%%%s%%", value)
	}

	return fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")), values
}

func profilesCsvSortings(params *list_params.ListParams) {
	params.AllowSortings([]string{"username", "email", "company_name", "first_name", "last_name", "created_at"})
}

func userFilter(inputValues []string, params *list_params.ListParams) (string, interface{}) {
	value := `%` + inputValues[0] + `%`
	return "(users.email LIKE ?)", []string{value}
}
