package basic

import (
	"fmt"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	permissionAdministrator "gitlab.com/iotTracker/brain/security/permission/administrator"
	permissionAdministratorException "gitlab.com/iotTracker/brain/security/permission/administrator/exception"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
	roleRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler"
	humanUser "gitlab.com/iotTracker/brain/user/human"
	userRecordHandler "gitlab.com/iotTracker/brain/user/human/recordHandler"
)

type administrator struct {
	userRecordHandler userRecordHandler.RecordHandler
	roleRecordHandler roleRecordHandler.RecordHandler
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
	roleRecordHandler roleRecordHandler.RecordHandler,
) permissionAdministrator.Administrator {
	return &administrator{
		userRecordHandler: userRecordHandler,
		roleRecordHandler: roleRecordHandler,
	}
}

func (a *administrator) ValidateUserHasPermissionRequest(request *permissionAdministrator.UserHasPermissionRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !humanUser.IsValidIdentifier(request.UserIdentifier) {
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

func (a *administrator) UserHasPermission(request *permissionAdministrator.UserHasPermissionRequest) (*permissionAdministrator.UserHasPermissionResponse, error) {
	if err := a.ValidateUserHasPermissionRequest(request); err != nil {
		return nil, err
	}

	// retrieve all of the users permissions
	getAllUsersPermissionsResponse, err := a.GetAllUsersAPIPermissions(&permissionAdministrator.GetAllUsersAPIPermissionsRequest{
		Claims:         request.Claims,
		UserIdentifier: request.UserIdentifier,
	})
	if err != nil {
		return nil, permissionAdministratorException.GetAllPermissions{Reasons: []string{err.Error()}}
	}

	response := permissionAdministrator.UserHasPermissionResponse{}

	// assume user does not have permission
	response.Result = false

	// go through all of the users permissions to see if one matches
	for _, perm := range getAllUsersPermissionsResponse.Permissions {
		if perm == request.Permission {
			response.Result = true
			return &response, nil
		}
	}

	return &response, nil
}

func (a *administrator) ValidateGetAllUsersAPIPermissionsRequest(request *permissionAdministrator.GetAllUsersAPIPermissionsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !humanUser.IsValidIdentifier(request.UserIdentifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.UserIdentifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (a *administrator) GetAllUsersAPIPermissions(request *permissionAdministrator.GetAllUsersAPIPermissionsRequest) (*permissionAdministrator.GetAllUsersAPIPermissionsResponse, error) {
	if err := a.ValidateGetAllUsersAPIPermissionsRequest(request); err != nil {
		return nil, err
	}

	// try and retrieve the user
	userRetrieveResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		return nil, err
	}

	usersAPIPermissions := make([]api.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range userRetrieveResponse.User.Roles {
		// retrieve the role
		roleRetrieveResponse, err := a.roleRecordHandler.Retrieve(&roleRecordHandler.RetrieveRequest{
			Identifier: name.Identifier{Name: roleName},
		})
		if err != nil {
			return nil, err
		}
		// add all of the permissions of the role
		usersAPIPermissions = append(usersAPIPermissions, roleRetrieveResponse.Role.APIPermissions...)
	}

	return &permissionAdministrator.GetAllUsersAPIPermissionsResponse{Permissions: usersAPIPermissions}, nil
}

func (a *administrator) ValidateGetAllUsersViewPermissionsRequest(request *permissionAdministrator.GetAllUsersViewPermissionsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !humanUser.IsValidIdentifier(request.UserIdentifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.UserIdentifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (a *administrator) GetAllUsersViewPermissions(request *permissionAdministrator.GetAllUsersViewPermissionsRequest) (*permissionAdministrator.GetAllUsersViewPermissionsResponse, error) {
	if err := a.ValidateGetAllUsersViewPermissionsRequest(request); err != nil {
		return nil, err
	}

	// try and retrieve the user
	userRetrieveResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.UserIdentifier,
	})
	if err != nil {
		return nil, err
	}

	usersViewPermissions := make([]view.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range userRetrieveResponse.User.Roles {
		// retrieve the role
		roleRetrieveResponse, err := a.roleRecordHandler.Retrieve(&roleRecordHandler.RetrieveRequest{
			Identifier: name.Identifier{Name: roleName},
		})
		if err != nil {
			return nil, err
		}
		// add all of the permissions of the role
		usersViewPermissions = append(usersViewPermissions, roleRetrieveResponse.Role.ViewPermissions...)
	}

	return &permissionAdministrator.GetAllUsersViewPermissionsResponse{Permissions: usersViewPermissions}, nil
}
