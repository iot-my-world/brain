package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/security/permission"
	permissionException "gitlab.com/iotTracker/brain/security/permission/exception"
	permissionHandler "gitlab.com/iotTracker/brain/security/permission/handler"
	roleRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler"
)

type handler struct {
	userRecordHandler userRecordHandler.RecordHandler
	roleRecordHandler roleRecordHandler.RecordHandler
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
	roleRecordHandler roleRecordHandler.RecordHandler,
) *handler {
	return &handler{
		userRecordHandler: userRecordHandler,
		roleRecordHandler: roleRecordHandler,
	}
}

func (bh *handler) ValidateUserHasPermissionRequest(request *permissionHandler.UserHasPermissionRequest) error {
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
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (bh *handler) UserHasPermission(request *permissionHandler.UserHasPermissionRequest, response *permissionHandler.UserHasPermissionResponse) error {
	if err := bh.ValidateUserHasPermissionRequest(request); err != nil {
		return err
	}

	// retrieve all of the users permissions
	getAllUsersPermissionsResponse := permissionHandler.GetAllUsersPermissionsResponse{}
	if err := bh.GetAllUsersPermissions(&permissionHandler.GetAllUsersPermissionsRequest{
		UserIdentifier: request.UserIdentifier,
	},
		&getAllUsersPermissionsResponse); err != nil {
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

func (bh *handler) ValidateGetAllUsersPermissionsRequest(request *permissionHandler.GetAllUsersPermissionsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.UserIdentifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.UserIdentifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (bh *handler) GetAllUsersPermissions(request *permissionHandler.GetAllUsersPermissionsRequest, response *permissionHandler.GetAllUsersPermissionsResponse) error {
	if err := bh.ValidateGetAllUsersPermissionsRequest(request); err != nil {
		return err
	}

	// try and retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := bh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{Identifier: request.UserIdentifier}, &userRetrieveResponse); err != nil {
		return err
	}

	usersPermissions := make([]permission.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range userRetrieveResponse.User.Roles {
		// retrieve the role
		roleRetrieveResponse := roleRecordHandler.RetrieveResponse{}
		if err := bh.roleRecordHandler.Retrieve(&roleRecordHandler.RetrieveRequest{Identifier: name.Identifier{Name: roleName}}, &roleRetrieveResponse); err != nil {
			return err
		}
		// add all of the permissions of the role
		usersPermissions = append(usersPermissions, roleRetrieveResponse.Role.Permissions...)
	}

	// return all permissions
	response.Permissions = usersPermissions

	return nil
}
