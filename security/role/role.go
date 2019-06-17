package role

import (
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/permission/view"
)

type Role struct {
	Id              string            `json:"id" bson:"id"`
	Name            string            `json:"name" bson:"name"`
	APIPermissions  []api.Permission  `json:"apiPermissions" bson:"apiPermissions"`
	ViewPermissions []view.Permission `json:"viewPermissions" bson:"viewPermissions"`
}

func (r *Role) CompareAPIPermissions(perms []api.Permission) bool {
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

func (r *Role) CompareViewPermissions(perms []view.Permission) bool {
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
