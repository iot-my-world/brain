package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/api"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	administrator2 "github.com/iot-my-world/brain/pkg/security/permission/administrator"
	"github.com/iot-my-world/brain/pkg/security/permission/administrator/exception"
	api2 "github.com/iot-my-world/brain/pkg/security/permission/api"
	view2 "github.com/iot-my-world/brain/pkg/security/permission/view"
	"github.com/iot-my-world/brain/pkg/security/role/recordHandler"
	apiUserRecordHandler "github.com/iot-my-world/brain/pkg/user/api/recordHandler"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	userRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
)

type administrator struct {
	userRecordHandler    userRecordHandler.RecordHandler
	roleRecordHandler    recordHandler.RecordHandler
	apiUserRecordHandler apiUserRecordHandler.RecordHandler
}

func New(
	userRecordHandler userRecordHandler.RecordHandler,
	roleRecordHandler recordHandler.RecordHandler,
	apiUserRecordHandler apiUserRecordHandler.RecordHandler,
) administrator2.Administrator {
	return &administrator{
		userRecordHandler:    userRecordHandler,
		roleRecordHandler:    roleRecordHandler,
		apiUserRecordHandler: apiUserRecordHandler,
	}
}

func (a *administrator) ValidateUserHasPermissionRequest(request *administrator2.UserHasPermissionRequest) error {
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

func (a *administrator) UserHasPermission(request *administrator2.UserHasPermissionRequest) (*administrator2.UserHasPermissionResponse, error) {
	if err := a.ValidateUserHasPermissionRequest(request); err != nil {
		return nil, err
	}

	// retrieve all of the users permissions
	getAllUsersPermissionsResponse, err := a.GetAllUsersAPIPermissions(&administrator2.GetAllUsersAPIPermissionsRequest{
		Claims:         request.Claims,
		UserIdentifier: request.UserIdentifier,
	})
	if err != nil {
		return nil, exception.GetAllPermissions{Reasons: []string{err.Error()}}
	}

	response := administrator2.UserHasPermissionResponse{}

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

func (a *administrator) ValidateGetAllUsersAPIPermissionsRequest(request *administrator2.GetAllUsersAPIPermissionsRequest) error {
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

func (a *administrator) GetAllUsersAPIPermissions(request *administrator2.GetAllUsersAPIPermissionsRequest) (*administrator2.GetAllUsersAPIPermissionsResponse, error) {
	if err := a.ValidateGetAllUsersAPIPermissionsRequest(request); err != nil {
		return nil, err
	}

	// get all of the roles assigned to this user
	var roles []string
	switch request.Claims.(type) {
	case human.Login:
		// try and retrieve the human user
		userRetrieveResponse, err := a.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.UserIdentifier,
		})
		if err != nil {
			return nil, err
		}
		roles = userRetrieveResponse.User.Roles

	case api.Login:
		// try and retrieve the api user
		apiUserRetrieveResponse, err := a.apiUserRecordHandler.Retrieve(&apiUserRecordHandler.RetrieveRequest{
			Claims:     request.Claims,
			Identifier: request.UserIdentifier,
		})
		if err != nil {
			return nil, err
		}
		roles = apiUserRetrieveResponse.User.Roles

	default:
		return nil, exception.GetAllPermissions{Reasons: []string{"invalid claims type", string(request.Claims.Type())}}
	}

	usersAPIPermissions := make([]api2.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range roles {
		// retrieve the role
		roleRetrieveResponse, err := a.roleRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
			Identifier: name.Identifier{Name: roleName},
		})
		if err != nil {
			return nil, err
		}
		// add all of the permissions of the role
		usersAPIPermissions = append(usersAPIPermissions, roleRetrieveResponse.Role.APIPermissions...)
	}

	return &administrator2.GetAllUsersAPIPermissionsResponse{Permissions: usersAPIPermissions}, nil
}

func (a *administrator) ValidateGetAllUsersViewPermissionsRequest(request *administrator2.GetAllUsersViewPermissionsRequest) error {
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

func (a *administrator) GetAllUsersViewPermissions(request *administrator2.GetAllUsersViewPermissionsRequest) (*administrator2.GetAllUsersViewPermissionsResponse, error) {
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

	usersViewPermissions := make([]view2.Permission, 0)

	// for every role that the user has been assigned
	for _, roleName := range userRetrieveResponse.User.Roles {
		// retrieve the role
		roleRetrieveResponse, err := a.roleRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
			Identifier: name.Identifier{Name: roleName},
		})
		if err != nil {
			return nil, err
		}
		// add all of the permissions of the role
		usersViewPermissions = append(usersViewPermissions, roleRetrieveResponse.Role.ViewPermissions...)
	}

	return &administrator2.GetAllUsersViewPermissionsResponse{Permissions: usersViewPermissions}, nil
}
