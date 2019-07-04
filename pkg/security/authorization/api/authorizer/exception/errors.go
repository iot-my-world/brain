package exception

import (
	"fmt"
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
)

type NotAuthorised struct {
	Permission api2.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}
