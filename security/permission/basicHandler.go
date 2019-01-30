package permission

import (
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/security/role"
	"fmt"
	globalException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/security"
	"gitlab.com/iotTracker/brain/search/identifiers/name"
	permissionException "gitlab.com/iotTracker/brain/security/permission/exception"
)

type basicHandler struct {
	userRecordHandler user.RecordHandler
	roleRecordHandler role.RecordHandler
}

func NewBasicHandler(
	userRecordHandler user.RecordHandler,
	roleRecordHandler role.RecordHandler,
) *basicHandler {
	return &basicHandler{
		userRecordHandler: userRecordHandler,
		roleRecordHandler: roleRecordHandler,
	}
}

func (bh *basicHandler) ValidateUserHasPermissionRequest(request *UserHasPermissionRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.UserIdentifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.UserIdentifier.Type()))
		}
	}

	if request.Permission == "" {
		reasonsInvalid = append(reasonsInvalid, "permission is blank")
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (bh *basicHandler) UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error {
	if err := bh.ValidateUserHasPermissionRequest(request); err != nil {
		return err
	}

	// retrieve all of the users permissions
	getAllUsersPermissionsResponse := GetAllUsersPermissionsResponse{}
	if err := bh.GetAllUsersPermissions(&GetAllUsersPermissionsRequest{}, &getAllUsersPermissionsResponse);
		err != nil {
		return permissionException.GetAllPermissions{Reasons: []string{err.Error()}}
	}

	// assume user does not have permission
	response.Result = false

	// go through all of the users permissions to see if one matches
	for _, perm := range getAllUsersPermissionsResponse.Permissions {
		if perm == request.Permission {
			response.Result = true
			return nil
		}
	}

	return nil
}

func (bh *basicHandler) ValidateGetAllUsersPermissionsRequest(request *GetAllUsersPermissionsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.UserIdentifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.UserIdentifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (bh *basicHandler) GetAllUsersPermissions(request *GetAllUsersPermissionsRequest, response *GetAllUsersPermissionsResponse) error {
	if err := bh.ValidateGetAllUsersPermissionsRequest(request); err != nil {
		return err
	}

	// try and retrieve the user
	userRetrieveResponse := user.RetrieveResponse{}
	if err := bh.userRecordHandler.Retrieve(&user.RetrieveRequest{Identifier: request.UserIdentifier}, &userRetrieveResponse);
		err != nil {
		return err
	}

	usersPermissions := make([]security.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range userRetrieveResponse.User.Roles {
		// retrieve the role
		roleRetrieveResponse := role.RetrieveResponse{}
		if err := bh.roleRecordHandler.Retrieve(&role.RetrieveRequest{Identifier: name.Identifier(roleName)}, &roleRetrieveResponse);
			err != nil {
			return err
		}
		// add all of the permissions of the role
		usersPermissions = append(usersPermissions, roleRetrieveResponse.Role.Permissions...)
	}

	// return all permissions
	response.Permissions = usersPermissions

	return nil
}
