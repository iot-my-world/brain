package apiAuth

import (
	"bitbucket.org/gotimekeeper/security/token"
	"bitbucket.org/gotimekeeper/security/systemRole"
	"bitbucket.org/gotimekeeper/security"
	"bitbucket.org/gotimekeeper/log"
	"errors"
)

type APIAuthorizer struct {
	JWTValidator token.JWTValidator
	RoleRecordHandler systemRole.RecordHandler
}

func (a *APIAuthorizer) AuthorizeAPIReq(jwt string, jsonRpcMethod string) error {

	// Validate the jwt
	jwtClaims, err := a.JWTValidator.ValidateJWT(jwt)
	if err != nil {
		return err
	}

	// TODO: Decide if we want to confirm that the user who owns this token has the role claimed in the token? i.e. Can we trust the signing process?

	// Check the if the user is authorised to access this jsonRpcMethod based on their role claim
	retrieveRoleResponse := systemRole.RetrieveResponse{}
	if err := a.RoleRecordHandler.Retrieve(&systemRole.RetrieveRequest{Name:jwtClaims.SystemRole}, &retrieveRoleResponse); err != nil {
		log.Info("Unable to retrieve role during API Access Authorisation!", err)
		return err
	}

	// Check if the jsonRpcMethod being accessed is the set of permissions assigned to their role
	for _, perm := range retrieveRoleResponse.SystemRole.Permissions {
		if perm == security.Permission(jsonRpcMethod) {
			return nil
		}
	}

	return errors.New("user does not have permission to access this API")
}
