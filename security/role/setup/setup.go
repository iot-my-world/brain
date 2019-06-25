package setup

import (
	"github.com/iot-my-world/brain/log"
	partyAdministrator "github.com/iot-my-world/brain/party/administrator"
	clientAdministrator "github.com/iot-my-world/brain/party/client/administrator"
	"github.com/iot-my-world/brain/search/identifier/id"
	"github.com/iot-my-world/brain/search/identifier/name"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/permission/view"
	"github.com/iot-my-world/brain/security/role"
	roleRecordHandler "github.com/iot-my-world/brain/security/role/recordHandler"
	roleRecordHandlerException "github.com/iot-my-world/brain/security/role/recordHandler/exception"
	roleSetupException "github.com/iot-my-world/brain/security/role/setup/exception"
)

var initialRoles = func() []role.Role {

	// Add Permissions to Roles
	// Party Administrator
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, partyAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, partyAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, partyAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, partyAdministrator.ClientUserPermissions...)

	// Client Administrator
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, clientAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, clientAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, clientAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, clientAdministrator.ClientUserPermissions...)

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
		// Role
		api.RoleCreate,
		api.RoleRetrieve,
		api.RoleUpdate,
		api.RoleDelete,

		// API User
		api.APIUserRecordHandlerCollect,
		api.APIUserAdministratorCreate,
		api.APIUserValidatorValidate,

		api.CompanyRecordHandlerRetrieve,

		api.CompanyValidatorValidate,

		api.CompanyAdministratorCreate,

		api.SystemAdministratorUpdateAllowedFields,

		api.PartyRegistrarInviteCompanyAdminUser,
		api.PartyRegistrarRegisterCompanyAdminUser,

		// Barcode Scanner
		api.BarcodeScannerScan,

		// SF001 Tracker
		api.SF001TrackerValidatorValidate,
		api.SF001TrackerAdministratorCreate,
		api.SF001TrackerAdministratorUpdateAllowedFields,
	}

	// The view permissions that root has
	rootViewPermissions := []view.Permission{
		view.PartyCompany,
		view.PartyClient,
		view.PartyUser,
		view.PartyAPIUser,

		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
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

var CompanyAdmin = role.Role{
	Name: "companyAdmin",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,

		api.SystemRecordHandlerCollect,

		// user
		api.UserRecordHandlerCollect,

		api.UserValidatorValidate,

		api.UserAdministratorUpdateAllowedFields,
		api.UserAdministratorCreate,
		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdatePassword,
		api.UserAdministratorCheckPassword,

		// Company
		api.CompanyRecordHandlerCollect,

		api.CompanyAdministratorUpdateAllowedFields,
		api.CompanyValidatorValidate,

		// Client
		api.ClientRecordHandlerRetrieve,
		api.ClientRecordHandlerCollect,

		api.ClientValidatorValidate,

		// Party
		api.PartyRegistrarInviteUser,
		api.PartyRegistrarInviteClientAdminUser,
		api.PartyRegistrarAreAdminsRegistered,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		// SF001 Tracker
		api.SF001TrackerRecordHandlerCollect,
	},
	ViewPermissions: []view.Permission{
		view.PartyProfileEditing,

		view.PartyClient,
		view.PartyUser,

		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
	},
}

var CompanyUser = role.Role{
	Name: "companyUser",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdateAllowedFields,
		api.UserAdministratorUpdatePassword,
		api.UserAdministratorCheckPassword,

		// SF001 Tracker
		api.SF001TrackerRecordHandlerCollect,
	},
	ViewPermissions: []view.Permission{
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
	},
}

var ClientAdmin = role.Role{
	Name: "clientAdmin",
	APIPermissions: []api.Permission{
		// user
		api.UserRecordHandlerCollect,

		api.UserValidatorValidate,

		api.UserAdministratorUpdateAllowedFields,
		api.UserAdministratorCreate,
		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdatePassword,
		api.UserAdministratorCheckPassword,

		// system
		api.SystemRecordHandlerCollect,

		// company
		api.CompanyRecordHandlerCollect,

		// client
		api.ClientRecordHandlerCollect,

		api.ClientValidatorValidate,

		api.PermissionHandlerGetAllUsersViewPermissions,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		api.PartyRegistrarInviteUser,
		api.PartyRegistrarAreAdminsRegistered,

		// SF001 Tracker
		api.SF001TrackerRecordHandlerCollect,
	},
	ViewPermissions: []view.Permission{
		view.PartyProfileEditing,

		view.PartyUser,

		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
	},
}
var ClientUser = role.Role{
	Name: "clientUser",
	APIPermissions: []api.Permission{
		api.PermissionHandlerGetAllUsersViewPermissions,

		api.TrackingReportLive,
		api.TrackingReportHistorical,

		api.UserAdministratorGetMyUser,
		api.UserAdministratorUpdateAllowedFields,
		api.UserAdministratorUpdatePassword,
		api.UserAdministratorCheckPassword,

		// SF001 Tracker
		api.SF001TrackerRecordHandlerCollect,
	},
	ViewPermissions: []view.Permission{
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
	},
}

func InitialSetup(handler roleRecordHandler.RecordHandler) error {
	for _, roleToCreate := range initialRoles {
		//Try and retrieve the record
		retrieveRoleResponse, err := handler.Retrieve(&roleRecordHandler.RetrieveRequest{
			Identifier: name.Identifier{Name: roleToCreate.Name},
		})
		switch err.(type) {
		case roleRecordHandlerException.NotFound:
			// if role record does not exist yet, try and create it
			_, err := handler.Create(&roleRecordHandler.CreateRequest{
				Role: roleToCreate,
			})
			if err != nil {
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
				if _, err := handler.Update(&roleRecordHandler.UpdateRequest{
					Role:       roleToCreate,
					Identifier: id.Identifier{Id: retrieveRoleResponse.Role.Id},
				}); err != nil {
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
