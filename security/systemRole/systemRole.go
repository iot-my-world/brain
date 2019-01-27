package systemRole

import (
	"gitlab.com/iotTracker/brain/security"
)

type SystemRole struct{
	Name string `json:"name" bson:"name"`
	Permissions []security.Permission `json:"permissions" bson:"permissions"`
}

func (r *SystemRole) comparePermissions(perms []security.Permission) bool {
	if len(r.Permissions) != len(perms) {
		return false
	}

	for i := range r.Permissions {
		if r.Permissions[i] != perms[i] {
		return false
		}
	}

	return true
}
