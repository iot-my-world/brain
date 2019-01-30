package exception

import (
	"gitlab.com/iotTracker/brain/security"
	"fmt"
)

type NotAuthorised struct {
	Permission security.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}