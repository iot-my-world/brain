package setup

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/security/permission/api"
	"gitlab.com/iotTracker/brain/security/permission/view"
	"gitlab.com/iotTracker/brain/security/role"
	roleRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler"
	roleRecordHandlerException "gitlab.com/iotTracker/brain/security/role/recordHandler/exception"
	roleSetupException "gitlab.com/iotTracker/brain/security/role/setup/exception"
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

		api.CompanyRecordHandlerRetrieve,

		api.CompanyValidatorValidate,

		api.CompanyAdministratorCreate,

		api.SystemAdministratorUpdateAllowedFields,

		api.PartyRegistrarInviteCompanyAdminUser,
		api.PartyRegistrarRegisterCompanyAdminUser,
		api.TK102DeviceAdministratorChangeOwnershipAndAssignment,
	}

	// The view permissions that root has
	rootViewPermissions := []view.Permission{
		view.Configuration,
		view.PartyCompanyConfiguration,
		view.PartyClientConfiguration,
		view.PartyUserConfiguration,
		view.DeviceConfiguration,
		view.Dashboards,
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,
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

		api.SystemRecordHandlerCollect,

		// user
		api.UserRecordHandlerCollect,

		api.UserValidatorValidate,

		api.UserAdministratorUpdateAllowedFields,

		// Company
		api.CompanyRecordHandlerCollect,

		api.CompanyAdministratorUpdateAllowedFields,

		// Client
		api.ClientRecordHandlerRetrieve,
		api.ClientRecordHandlerCollect,

		api.ClientValidatorValidate,

		api.ClientAdministratorCreate,
		api.ClientAdministratorUpdateAllowedFields,

		// Party
		api.PartyRegistrarInviteCompanyUser,
		api.PartyRegistrarInviteClientAdminUser,
		api.PartyRegistrarAreAdminsRegistered,

		// TK102 Device
		api.TK102DeviceRecordHandlerCreate,
		api.TK102DeviceRecordHandlerRetrieve,
		api.TK102DeviceRecordHandlerDelete,
		api.TK102DeviceRecordHandlerValidate,
		api.TK102DeviceRecordHandlerCollect,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		api.ReadingRecordHandlerCollect,

		api.PartyAdministratorGetMyParty,

		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdateAllowedFields,
	},
	ViewPermissions: []view.Permission{
		view.Configuration,
		view.PartyClientConfiguration,
		view.PartyUserConfiguration,
		view.DeviceConfiguration,
		view.Dashboards,
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,
	},
}
var CompanyUser = role.Role{
	Name: "companyUser",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		api.ReadingRecordHandlerCollect,

		api.PartyAdministratorGetMyParty,

		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdateAllowedFields,
	},
	ViewPermissions: []view.Permission{
		view.Dashboards,
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,
	},
}

var ClientAdmin = role.Role{
	Name: "clientAdmin",
	APIPermissions: []api.Permission{
		// user
		api.UserRecordHandlerCollect,

		api.UserValidatorValidate,

		api.UserAdministratorUpdateAllowedFields,

		// system
		api.SystemRecordHandlerCollect,

		// company
		api.CompanyRecordHandlerCollect,

		// client
		api.ClientRecordHandlerCollect,

		api.ClientAdministratorUpdateAllowedFields,

		api.PermissionHandlerGetAllUsersViewPermissions,

		api.TK102DeviceRecordHandlerCreate,
		api.TK102DeviceRecordHandlerRetrieve,
		api.TK102DeviceRecordHandlerDelete,
		api.TK102DeviceRecordHandlerValidate,
		api.TK102DeviceRecordHandlerCollect,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		api.ReadingRecordHandlerCollect,

		api.PartyRegistrarInviteClientUser,
		api.PartyRegistrarAreAdminsRegistered,

		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdateAllowedFields,

		api.PartyAdministratorGetMyParty,
	},
	ViewPermissions: []view.Permission{
		view.Configuration,
		view.PartyClientConfiguration,
		view.PartyUserConfiguration,
		view.Dashboards,
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,
	},
}

var ClientUser = role.Role{
	Name: "clientUser",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		api.ReadingRecordHandlerCollect,

		api.PartyAdministratorGetMyParty,

		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdateAllowedFields,
	},
	ViewPermissions: []view.Permission{
		view.Dashboards,
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,
	},
}

func InitialSetup(handler roleRecordHandler.RecordHandler) error {
	for _, roleToCreate := range initialRoles {
		//Try and retrieve the record
		retrieveRoleResponse := roleRecordHandler.RetrieveResponse{}
		err := handler.Retrieve(&roleRecordHandler.RetrieveRequest{Identifier: name.Identifier{Name: roleToCreate.Name}}, &retrieveRoleResponse)

		switch err.(type) {
		case roleRecordHandlerException.NotFound:
			// if role record does not exist yet, try and create it
			createRoleResponse := roleRecordHandler.CreateResponse{}
			if err := handler.Create(&roleRecordHandler.CreateRequest{Role: roleToCreate}, &createRoleResponse); err != nil {
				return roleSetupException.InitialSetup{Reasons: []string{"creation error", err.Error()}}
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
					return roleSetupException.InitialSetup{Reasons: []string{"update error", err.Error()}}
				}
			}

		default:
			// otherwise there was some retrieval error
			return roleSetupException.InitialSetup{Reasons: []string{"retrieval error", err.Error()}}
		}
	}

	return nil
}
