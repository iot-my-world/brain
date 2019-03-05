package administrator

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/tracker/device/tk102"
)

type Administrator interface {
	ChangeOwner(request *ChangeOwnerRequest, response *ChangeOwnerResponse) error
	ChangeAssigned(request *ChangeAssignedRequest, response *ChangeAssignedResponse) error
}

type ChangeOwnerRequest struct {
	Claims             claims.Claims
	TK02Identifier     identifier.Identifier
	NewOwnerPartyType  party.Type
	NewOwnerIdentifier identifier.Identifier
}

type ChangeOwnerResponse struct {
	TK102 tk102.TK102
}

type ChangeAssignedRequest struct {
	Claims                claims.Claims
	TK02Identifier        identifier.Identifier
	NewAssignedPartyType  party.Type
	NewAssignedIdentifier identifier.Identifier
}

type ChangeAssignedResponse struct {
	TK102 tk102.TK102
}
