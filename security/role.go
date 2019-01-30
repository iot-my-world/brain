package security

type Role struct {
	Id          string       `json:"id" bson:"id"`
	Name        string       `json:"name" bson:"name"`
	Permissions []Permission `json:"permissions" bson:"permissions"`
}

func (r *Role) ComparePermissions(perms []Permission) bool {
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
