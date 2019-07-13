package exception

import (
	"fmt"
	apiPermission "github.com/iot-my-world/brain/pkg/security/permission/api"
)

type NotAuthorised struct {
	Permission apiPermission.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}
