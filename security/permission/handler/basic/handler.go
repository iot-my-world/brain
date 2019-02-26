package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/security/permission/api"
	permissionHandlerException "gitlab.com/iotTracker/brain/security/permission/handler/exception"
	permissionHandler "gitlab.com/iotTracker/brain/security/permission/handler"
	"gitlab.com/iotTracker/brain/security/permission/view"
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
	getAllUsersPermissionsResponse := permissionHandler.GetAllUsersAPIPermissionsResponse{}
	if err := bh.GetAllUsersAPIPermissions(&permissionHandler.GetAllUsersAPIPermissionsRequest{
		UserIdentifier: request.UserIdentifier,
	},
		&getAllUsersPermissionsResponse); err != nil {
		return permissionHandlerException.GetAllPermissions{Reasons: []string{err.Error()}}
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

func (bh *handler) ValidateGetAllUsersAPIPermissionsRequest(request *permissionHandler.GetAllUsersAPIPermissionsRequest) error {
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

func (bh *handler) GetAllUsersAPIPermissions(request *permissionHandler.GetAllUsersAPIPermissionsRequest, response *permissionHandler.GetAllUsersAPIPermissionsResponse) error {
	if err := bh.ValidateGetAllUsersAPIPermissionsRequest(request); err != nil {
		return err
	}

	// try and retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := bh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{Identifier: request.UserIdentifier}, &userRetrieveResponse); err != nil {
		return err
	}

	usersAPIPermissions := make([]api.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range userRetrieveResponse.User.Roles {
		// retrieve the role
		roleRetrieveResponse := roleRecordHandler.RetrieveResponse{}
		if err := bh.roleRecordHandler.Retrieve(&roleRecordHandler.RetrieveRequest{Identifier: name.Identifier{Name: roleName}}, &roleRetrieveResponse); err != nil {
			return err
		}
		// add all of the permissions of the role
		usersAPIPermissions = append(usersAPIPermissions, roleRetrieveResponse.Role.APIPermissions...)
	}

	// return all permissions
	response.Permissions = usersAPIPermissions

	return nil
}

func (bh *handler) ValidateGetAllUsersViewPermissionsRequest(request *permissionHandler.GetAllUsersViewPermissionsRequest) error {
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

func (bh *handler) GetAllUsersViewPermissions(request *permissionHandler.GetAllUsersViewPermissionsRequest, response *permissionHandler.GetAllUsersViewPermissionsResponse) error {
	if err := bh.ValidateGetAllUsersViewPermissionsRequest(request); err != nil {
		return err
	}

	// try and retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := bh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{Identifier: request.UserIdentifier}, &userRetrieveResponse); err != nil {
		return err
	}

	usersViewPermissions := make([]view.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range userRetrieveResponse.User.Roles {
		// retrieve the role
		roleRetrieveResponse := roleRecordHandler.RetrieveResponse{}
		if err := bh.roleRecordHandler.Retrieve(&roleRecordHandler.RetrieveRequest{Identifier: name.Identifier{Name: roleName}}, &roleRetrieveResponse); err != nil {
			return err
		}
		// add all of the permissions of the role
		usersViewPermissions = append(usersViewPermissions, roleRetrieveResponse.Role.ViewPermissions...)
	}

	// return all permissions
	response.Permissions = usersViewPermissions

	return nil
}
