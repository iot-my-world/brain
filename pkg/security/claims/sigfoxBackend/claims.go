package sigfoxBackend

import (
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/security/claims"
	apiPermission "github.com/iot-my-world/brain/pkg/security/permission/api"
	sigfoxBackendCallbackServer "github.com/iot-my-world/brain/pkg/sigfox/backend/callback/server"
	"time"
)

type SigfoxBackend struct {
	BackendId id.Identifier `json:"backendId"`
}

func (r SigfoxBackend) Type() claims.Type {
	return claims.SigfoxBackend
}

func (r SigfoxBackend) Expired() bool {
	// these claims never expire
	return false
}

func (r SigfoxBackend) TimeToExpiry() time.Duration {
	return -1
}

func (r SigfoxBackend) PartyDetails() party.Details {
	return party.Details{}
}

// permissions granted by having a valid set of these claims
var GrantedAPIPermissions = []apiPermission.Permission{
	sigfoxBackendCallbackServer.HandleDataMessageService,
}
