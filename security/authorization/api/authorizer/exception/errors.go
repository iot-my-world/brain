package exception

import (
	"fmt"
	"gitlab.com/iotTracker/brain/security/permission/api"
)

type NotAuthorised struct {
	Permission api.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}
