package setup

import (
	"github.com/iot-my-world/brain/log"
	partyAdministrator "github.com/iot-my-world/brain/pkg/party/administrator"
	clientAdministrator "github.com/iot-my-world/brain/pkg/party/client/administrator"
	clientRecordHandler "github.com/iot-my-world/brain/pkg/party/client/recordHandler"
	clientValidator "github.com/iot-my-world/brain/pkg/party/client/validator"
	companyAdministrator "github.com/iot-my-world/brain/pkg/party/company/administrator"
	companyRecordHandler "github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	companyValidator "github.com/iot-my-world/brain/pkg/party/company/validator"
	partyRegistrar "github.com/iot-my-world/brain/pkg/party/registrar"
	systemRecordHandler "github.com/iot-my-world/brain/pkg/party/system/recordHandler"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	"github.com/iot-my-world/brain/pkg/search/identifier/name"
	sf001TrackerRecordHandler "github.com/iot-my-world/brain/pkg/tracker/sf001/recordHandler"
	sf001TrackerValidator "github.com/iot-my-world/brain/pkg/tracker/sf001/validator"
	humanUserAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
	humanUserRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	humanUserValidator "github.com/iot-my-world/brain/pkg/user/human/validator"
	permissionAdministrator "github.com/iot-my-world/brain/security/permission/administrator"
	"github.com/iot-my-world/brain/security/permission/api"
	"github.com/iot-my-world/brain/security/permission/view"
	"github.com/iot-my-world/brain/security/role"
	roleRecordHandler "github.com/iot-my-world/brain/security/role/recordHandler"
	roleRecordHandlerException "github.com/iot-my-world/brain/security/role/recordHandler/exception"
	roleSetupException "github.com/iot-my-world/brain/security/role/setup/exception"
)

var CompanyAdmin = role.Role{
	Name:           "companyAdmin",
	APIPermissions: make([]api.Permission, 0),
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
	Name:           "companyUser",
	APIPermissions: make([]api.Permission, 0),
	ViewPermissions: []view.Permission{
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
	},
}

var ClientAdmin = role.Role{
	Name:           "clientAdmin",
	APIPermissions: make([]api.Permission, 0),
	ViewPermissions: []view.Permission{
		view.PartyProfileEditing,

		view.PartyUser,

		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
	},
}

var ClientUser = role.Role{
	Name:           "clientUser",
	APIPermissions: make([]api.Permission, 0),
	ViewPermissions: []view.Permission{
		view.LiveTrackingDashboard,
		view.HistoricalTrackingDashboard,

		view.TrackerSF001,
	},
}

var initialRoles = func() []role.Role {

	rootAPIPermissions := make([]api.Permission, 0)

	// System RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, systemRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, systemRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, systemRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, systemRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, systemRecordHandler.ClientUserPermissions...)

	// Role RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, roleRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, roleRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, roleRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, roleRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, roleRecordHandler.ClientUserPermissions...)

	// Permission Administrator
	rootAPIPermissions = append(rootAPIPermissions, permissionAdministrator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, permissionAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, permissionAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, permissionAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, permissionAdministrator.ClientUserPermissions...)

	// Human User RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, humanUserRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, humanUserRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, humanUserRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, humanUserRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, humanUserRecordHandler.ClientUserPermissions...)
	// Human User Administrator
	rootAPIPermissions = append(rootAPIPermissions, humanUserAdministrator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, humanUserAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, humanUserAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, humanUserAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, humanUserAdministrator.ClientUserPermissions...)
	// Human User Validator
	rootAPIPermissions = append(rootAPIPermissions, humanUserValidator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, humanUserValidator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, humanUserValidator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, humanUserValidator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, humanUserValidator.ClientUserPermissions...)

	// Party Administrator
	rootAPIPermissions = append(rootAPIPermissions, partyAdministrator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, partyAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, partyAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, partyAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, partyAdministrator.ClientUserPermissions...)
	// Party Registrar
	rootAPIPermissions = append(rootAPIPermissions, partyRegistrar.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, partyRegistrar.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, partyRegistrar.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, partyRegistrar.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, partyRegistrar.ClientUserPermissions...)

	// Company Administrator
	rootAPIPermissions = append(rootAPIPermissions, companyAdministrator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, companyAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, companyAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, companyAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, companyAdministrator.ClientUserPermissions...)
	// Company RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, companyRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, companyRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, companyRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, companyRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, companyRecordHandler.ClientUserPermissions...)
	// Company Validator
	rootAPIPermissions = append(rootAPIPermissions, companyValidator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, companyValidator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, companyValidator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, companyValidator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, companyValidator.ClientUserPermissions...)

	// Client Administrator
	rootAPIPermissions = append(rootAPIPermissions, clientAdministrator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, clientAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, clientAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, clientAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, clientAdministrator.ClientUserPermissions...)
	// Client RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, clientRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, clientRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, clientRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, clientRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, clientRecordHandler.ClientUserPermissions...)
	// Client Validator
	rootAPIPermissions = append(rootAPIPermissions, clientValidator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, clientValidator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, clientValidator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, clientValidator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, clientValidator.ClientUserPermissions...)

	// SF001 Tracker RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, sf001TrackerRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, sf001TrackerRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, sf001TrackerRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, sf001TrackerRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, sf001TrackerRecordHandler.ClientUserPermissions...)

	// SF001 Tracker Validator
	rootAPIPermissions = append(rootAPIPermissions, sf001TrackerValidator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, sf001TrackerValidator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, sf001TrackerValidator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, sf001TrackerValidator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, sf001TrackerValidator.ClientUserPermissions...)

	// Register roles here
	allRoles := []role.Role{
		ClientAdmin,
		ClientUser,
		CompanyAdmin,
		CompanyUser,
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
