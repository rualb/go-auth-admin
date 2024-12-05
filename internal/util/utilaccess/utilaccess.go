package utilaccess

import (
	"slices"
	"strings"
)

const (
	RoleAdmin = "admin"
)

type PermissionsDTO struct {
	Access bool `json:"access,omitempty"`
	Add    bool `json:"add,omitempty"`
	View   bool `json:"view,omitempty"`
	Edit   bool `json:"edit,omitempty"`
	Delete bool `json:"delete,omitempty"`
}

func (x *PermissionsDTO) Fill(userRoles string, prefix string) {
	roles := strings.Fields(userRoles)

	for _, r := range roles {
		if IsAdmin(r) {
			x.Access = true
			x.Add = true
			x.View = true
			x.Edit = true
			x.Delete = true
			break
		}
		if !strings.HasPrefix(r, prefix) {
			continue
		}
		switch {
		case r == prefix+"access":
			x.Access = true
		case r == prefix+"add":
			x.Add = true
		case r == prefix+"view":
			x.View = true
		case r == prefix+"edit":
			x.Edit = true
		case r == prefix+"delete":
			x.Delete = true
		}
	}
}

func IsAdmin(role string) bool {
	return role == RoleAdmin
}

// // HasAllRoles checks if the user has all the specified roles.
// func HasAllRoles(userRoles string, roles ...string) bool {

// 	// If no roles are provided, return false
// 	if userRoles == "" || len(roles) == 0 {
// 		return false
// 	}

// 	// Split the user's roles into a slice
// 	userRolesArr := strings.Fields(userRoles)

// 	for _, r := range userRolesArr {
// 		if IsAdmin(r) {
// 			return true
// 		}
// 	}
// 	// Check if the user has all specified roles

// 	for _, r := range roles {
// 		hasRole := slices.Contains(userRolesArr, r)
// 		if !hasRole { // !!! "!hasRole"
// 			return false
// 		}
// 	}

// 	return true // has all roles
// }

// HasAnyOfRoles checks if the user has all the specified roles.
func HasAnyOfRoles(userRoles string, roles ...string) bool {

	// If no roles are provided, return false
	if userRoles == "" || len(roles) == 0 {
		return false
	}

	// Split the user's roles into a slice
	userRolesArr := strings.Fields(userRoles)

	for _, r := range userRolesArr {
		if IsAdmin(r) {
			return true
		}
	}
	// Check if the user has all specified roles

	for _, r := range roles {
		hasRole := slices.Contains(userRolesArr, r)
		if hasRole { // !!! "hasRole"
			return true
		}
	}

	return false // has all roles
}
