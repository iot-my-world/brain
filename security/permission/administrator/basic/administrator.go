package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	permissionAdministrator "gitlab.com/iotTracker/brain/security/permission/administrator"
	permissionAdministratorException "gitlab.com/iotTracker/brain/security/permission/administrator/exception"
	"gitlab.com/iotTracker/brain/security/permission/api"
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

func (bh *handler) ValidateUserHasPermissionRequest(request *permissionAdministrator.UserHasPermissionRequest) error {
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

func (bh *handler) UserHasPermission(request *permissionAdministrator.UserHasPermissionRequest, response *permissionAdministrator.UserHasPermissionResponse) error {
	if err := bh.ValidateUserHasPermissionRequest(request); err != nil {
		return err
	}

	// retrieve all of the users permissions
	getAllUsersPermissionsResponse := permissionAdministrator.GetAllUsersAPIPermissionsResponse{}
	if err := bh.GetAllUsersAPIPermissions(&permissionAdministrator.GetAllUsersAPIPermissionsRequest{
		Claims:         request.Claims,
		UserIdentifier: request.UserIdentifier,
	},
		&getAllUsersPermissionsResponse); err != nil {
		return permissionAdministratorException.GetAllPermissions{Reasons: []string{err.Error()}}
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

func (bh *handler) ValidateGetAllUsersAPIPermissionsRequest(request *permissionAdministrator.GetAllUsersAPIPermissionsRequest) error {
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

func (bh *handler) GetAllUsersAPIPermissions(request *permissionAdministrator.GetAllUsersAPIPermissionsRequest, response *permissionAdministrator.GetAllUsersAPIPermissionsResponse) error {
	if err := bh.ValidateGetAllUsersAPIPermissionsRequest(request); err != nil {
		return err
	}

	// try and retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := bh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	}, &userRetrieveResponse); err != nil {
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

func (bh *handler) ValidateGetAllUsersViewPermissionsRequest(request *permissionAdministrator.GetAllUsersViewPermissionsRequest) error {
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

func (bh *handler) GetAllUsersViewPermissions(request *permissionAdministrator.GetAllUsersViewPermissionsRequest, response *permissionAdministrator.GetAllUsersViewPermissionsResponse) error {
	if err := bh.ValidateGetAllUsersViewPermissionsRequest(request); err != nil {
		return err
	}

	// try and retrieve the user
	userRetrieveResponse := userRecordHandler.RetrieveResponse{}
	if err := bh.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	}, &userRetrieveResponse); err != nil {
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
