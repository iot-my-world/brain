package exception

import "strings"

type GetAllPermissions struct {
	Reasons []string
}

func (e GetAllPermissions) Error() string {
	return "error getting all users permissions: " + strings.Join(e.Reasons,  "; ")
}
