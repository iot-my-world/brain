package exception

import (
	"fmt"
	"github.com/iot-my-world/brain/pkg/security/claims"
	apiPermission "github.com/iot-my-world/brain/pkg/security/permission/api"
)

type NotAuthorised struct {
	Permission apiPermission.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}

type InvalidClaims struct {
	ExpectedClaimsType claims.Type
}

func (e InvalidClaims) Error() string {
	return fmt.Sprintf("invalid claims, expected: %s", e.ExpectedClaimsType)
}
