package setup

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
	"gitlab.com/iotTracker/brain/security/role"
	roleException "gitlab.com/iotTracker/brain/security/role/exception"
	roleRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler"
)

var initialRoles = func() []role.Role {

	// Register roles here
	allRoles := []role.Role{
		ClientAdmin,
		ClientUser,
		CompanyAdmin,
		CompanyUser,
	}

	// Register additional root api permissions here
	// i.e. these are permissions that ONLY root has
	rootAPIPermissions := []api.Permission{
		api.RoleCreate,
		api.RoleRetrieve,
		api.RoleUpdate,
		api.RoleDelete,
		api.CompanyRecordHandlerCreate,
		api.CompanyRecordHandlerRetrieve,
		api.CompanyRecordHandlerUpdate,
		api.CompanyRecordHandlerDelete,
		api.CompanyRecordHandlerValidate,
		api.CompanyRecordHandlerCollect,
		api.PartyRegistrarInviteCompanyAdminUser,
		api.PartyRegistrarRegisterCompanyAdminUser,
		api.ReadingRecordHandlerCollect,
		api.TK102DeviceAdministratorChangeOwner,
		api.TK102DeviceAdministratorChangeAssigned,
	}

	rootViewPermissions := []view.Permission{
		view.Configuration,
		view.PartyCompanyConfiguration,
		view.PartyClientConfiguration,
		view.PartyUserConfiguration,
		view.DeviceConfiguration,
	}

	// Create root role and apply permissions of all other roles to root
	for _, role := range allRoles {
		// for each api permission in this role
	RoleAPIPerms:
		for _, apiPerm := range role.APIPermissions {
			// check of root already has it
			for _, rootAPIPerm := range rootAPIPermissions {
				if rootAPIPerm == apiPerm {
					continue RoleAPIPerms
				}
			}
			// if we are here root doesn't have it yet
			rootAPIPermissions = append(rootAPIPermissions, apiPerm)
		}
	}
	root := role.Role{
		Name:            "root",
		APIPermissions:  rootAPIPermissions,
		ViewPermissions: rootViewPermissions,
	}
	return append([]role.Role{root}, allRoles...)
}()

// Create Roles here
var CompanyAdmin = role.Role{
	Name: "companyAdmin",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,
		api.ClientRecordHandlerCreate,
		api.ClientRecordHandlerRetrieve,
		api.ClientRecordHandlerUpdate,
		api.ClientRecordHandlerDelete,
		api.ClientRecordHandlerValidate,
		api.ClientRecordHandlerCollect,
		api.PartyRegistrarInviteClientAdminUser,
		api.TK102DeviceRecordHandlerCreate,
		api.TK102DeviceRecordHandlerRetrieve,
		api.TK102DeviceRecordHandlerDelete,
		api.TK102DeviceRecordHandlerValidate,
		api.TK102DeviceRecordHandlerCollect,
		api.TK102DeviceAdministratorChangeAssigned,
	},
	ViewPermissions: []view.Permission{
		view.Configuration,
		view.PartyClientConfiguration,
		view.PartyUserConfiguration,
		view.DeviceConfiguration,
	},
}
var CompanyUser = role.Role{
	Name: "companyUser",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,
	},
	ViewPermissions: []view.Permission{},
}

var ClientAdmin = role.Role{
	Name: "clientAdmin",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,
		api.TK102DeviceRecordHandlerCreate,
		api.TK102DeviceRecordHandlerRetrieve,
		api.TK102DeviceRecordHandlerDelete,
		api.TK102DeviceRecordHandlerValidate,
		api.TK102DeviceRecordHandlerCollect,
	},
	ViewPermissions: []view.Permission{
		view.Configuration,
		view.PartyClientConfiguration,
		view.PartyUserConfiguration,
	},
}

var ClientUser = role.Role{
	Name: "clientUser",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,
	},
	ViewPermissions: []view.Permission{},
}

func InitialSetup(handler roleRecordHandler.RecordHandler) error {
	for _, roleToCreate := range initialRoles {
		//Try and retrieve the record
		retrieveRoleResponse := roleRecordHandler.RetrieveResponse{}
		err := handler.Retrieve(&roleRecordHandler.RetrieveRequest{Identifier: name.Identifier{Name: roleToCreate.Name}}, &retrieveRoleResponse)

		switch err.(type) {
		case roleException.NotFound:
			// if role record does not exist yet, try and create it
			createRoleResponse := roleRecordHandler.CreateResponse{}
			if err := handler.Create(&roleRecordHandler.CreateRequest{Role: roleToCreate}, &createRoleResponse); err != nil {
				return roleException.InitialSetup{Reasons: []string{"creation error", err.Error()}}
			}
			log.Info("Initial Role Setup: Created Role: " + roleToCreate.Name)

		case nil:
			//Record Retrieved Successfully

			// Update Role Permissions If Necessary
			if !(roleToCreate.CompareAPIPermissions(retrieveRoleResponse.Role.APIPermissions) &&
				roleToCreate.CompareViewPermissions(retrieveRoleResponse.Role.ViewPermissions)) {
				// role permissions differ, try update role
				log.Info("Initial Role Setup: Role: " + roleToCreate.Name + " already exists. Updating Role API permissions.")
				if err := handler.Update(&roleRecordHandler.UpdateRequest{
					Role:       roleToCreate,
					Identifier: id.Identifier{Id: retrieveRoleResponse.Role.Id},
				},
					&roleRecordHandler.UpdateResponse{}); err != nil {
					return roleException.InitialSetup{Reasons: []string{"update error", err.Error()}}
				}
			}

		default:
			// otherwise there was some retrieval error
			return roleException.InitialSetup{Reasons: []string{"retrieval error", err.Error()}}
		}
	}

	return nil
}
