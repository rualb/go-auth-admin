package utilaccess

import (
	"reflect"
	"testing"
)

func TestPermissionsDTO(t *testing.T) {
	type args struct {
		userRoles string
		prefix    string
	}

	tests := []struct {
		name string
		x    *PermissionsDTO
		args args
		want *PermissionsDTO
	}{
		{
			name: "Empty roles",
			x:    &PermissionsDTO{},
			args: args{"", "user_"},
			want: &PermissionsDTO{},
		},
		{
			name: "No roles with prefix",
			x:    &PermissionsDTO{},
			args: args{"admin user_access", "user_"},
			want: &PermissionsDTO{Access: true, Add: true, View: true, Edit: true, Delete: true},
		},
		{
			name: "Single role with prefix - access",
			x:    &PermissionsDTO{},
			args: args{"user_access", "user_"},
			want: &PermissionsDTO{Access: true},
		},
		{
			name: "Multiple roles with prefix",
			x:    &PermissionsDTO{},
			args: args{"user_access user_add user_edit", "user_"},
			want: &PermissionsDTO{Access: true, Add: true, Edit: true},
		},
		{
			name: "Roles with mixed prefixes",
			x:    &PermissionsDTO{},
			args: args{"user_access user_delete", "user_"},
			want: &PermissionsDTO{Access: true, Delete: true},
		},
		{
			name: "All roles",
			x:    &PermissionsDTO{},
			args: args{"user_access user_add user_edit user_delete", "user_"},
			want: &PermissionsDTO{Access: true, Add: true, Edit: true, Delete: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.x.Fill(tt.args.userRoles, tt.args.prefix)
			if !reflect.DeepEqual(tt.x, tt.want) {
				t.Errorf("PermissionsDTO.Fill() = %v, want %v", tt.x, tt.want)
			}
		})
	}
}
