package exception

import (
	"fmt"
	"gitlab.com/iotTracker/brain/security/permission"
)

type NotAuthorised struct {
	Permission permission.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}
