package apiAuth

import (
	"gitlab.com/iotTracker/brain/security/token"
	"gitlab.com/iotTracker/brain/security/role"
	"gitlab.com/iotTracker/brain/security"
	"gitlab.com/iotTracker/brain/log"
	"errors"
)

type APIAuthorizer struct {
	JWTValidator token.JWTValidator
	RoleRecordHandler role.RecordHandler
}

func (a *APIAuthorizer) AuthorizeAPIReq(jwt string, jsonRpcMethod string) error {

	// Validate the jwt
	jwtClaims, err := a.JWTValidator.ValidateJWT(jwt)
	if err != nil {
		return err
	}

	// Check the if the user is authorised to access this jsonRpcMethod based on their role claim
	retrieveRoleResponse := role.RetrieveResponse{}
	if err := a.RoleRecordHandler.Retrieve(&role.RetrieveRequest{Name:jwtClaims.Role}, &retrieveRoleResponse); err != nil {
		log.Info("Unable to retrieve role during API Access Authorisation!", err)
		return err
	}

	// Check if the jsonRpcMethod being accessed is the set of permissions assigned to their role
	for _, perm := range retrieveRoleResponse.Role.Permissions {
		if perm == security.Permission(jsonRpcMethod) {
			return nil
		}
	}

	return errors.New("user does not have permission to access this API")
}
