package authadmin

import (
	"go-auth-admin/internal/config/consts"
	"net/http"

	"github.com/labstack/echo/v4"
)

var mapPathToRoles = map[string][]string{
	// POST /api HTTP/1.1

	http.MethodGet + " " + consts.PathAuthAdminStatusAPI: {consts.AuthRoleAccess},
	http.MethodGet + " " + consts.PathAuthAdminConfigAPI: {consts.AuthRoleAccess},

	http.MethodGet + " " + consts.PathAuthAdminAccounts:       {consts.AuthRoleAccess},
	http.MethodGet + " " + consts.PathAuthAdminAccountsEntity: {consts.AuthRoleAccess},
	http.MethodGet + " " + consts.PathAuthAdminAccountsAPI:    {consts.AuthRoleAccess},

	http.MethodGet + " " + consts.PathAuthAdminAccountsEntityAPI:    {consts.AuthRoleEdit, consts.AuthRoleView},
	http.MethodPost + " " + consts.PathAuthAdminAccountsAPI:         {consts.AuthRoleAdd},
	http.MethodPut + " " + consts.PathAuthAdminAccountsEntityAPI:    {consts.AuthRoleEdit},
	http.MethodDelete + " " + consts.PathAuthAdminAccountsEntityAPI: {consts.AuthRoleDelete},

	http.MethodGet + " " + consts.PathAuthAdminAccountsEntityByCodeAPI: {consts.AuthRoleAccess},

	http.MethodPost + " " + consts.PathAuthAdminAccountsEntityPasswordAPI: {consts.AuthRoleEdit},
}

func RolesForAPI(c echo.Context) []string {
	key := c.Request().Method + " " + c.Path() // Method + Grop path + Route path
	roles := mapPathToRoles[key]
	return roles
}

func RolesForAssets(_ echo.Context) []string {
	return []string{consts.AuthRoleAccess}
}
