package role

import (
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
	view2 "github.com/iot-my-world/brain/pkg/security/permission/view"
)

type Role struct {
	Id              string             `json:"id" bson:"id"`
	Name            string             `json:"name" bson:"name"`
	APIPermissions  []api2.Permission  `json:"apiPermissions" bson:"apiPermissions"`
	ViewPermissions []view2.Permission `json:"viewPermissions" bson:"viewPermissions"`
}

func (r *Role) CompareAPIPermissions(perms []api2.Permission) bool {
	if len(r.APIPermissions) != len(perms) {
		return false
	}

	for i := range r.APIPermissions {
		if r.APIPermissions[i] != perms[i] {
			return false
		}
	}

	return true
}

func (r *Role) CompareViewPermissions(perms []view2.Permission) bool {
	if len(r.ViewPermissions) != len(perms) {
		return false
	}

	for i := range r.ViewPermissions {
		if r.ViewPermissions[i] != perms[i] {
			return false
		}
	}

	return true
}
