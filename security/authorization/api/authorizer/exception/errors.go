package exception

import (
	"fmt"
	"github.com/iot-my-world/brain/security/permission/api"
)

type NotAuthorised struct {
	Permission api.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}
