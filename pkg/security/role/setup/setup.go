package setup

import (
	"github.com/iot-my-world/brain/internal/log"
	sigbugAdministrator "github.com/iot-my-world/brain/pkg/device/sigbug/administrator"
	sigbugRecordHandler "github.com/iot-my-world/brain/pkg/device/sigbug/recordHandler"
	sigbugValidator "github.com/iot-my-world/brain/pkg/device/sigbug/validator"
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
	"github.com/iot-my-world/brain/pkg/security/permission/administrator"
	apiPermission "github.com/iot-my-world/brain/pkg/security/permission/api"
	viewPermission "github.com/iot-my-world/brain/pkg/security/permission/view"
	"github.com/iot-my-world/brain/pkg/security/role"
	"github.com/iot-my-world/brain/pkg/security/role/recordHandler"
	roleRecordHandlerException "github.com/iot-my-world/brain/pkg/security/role/recordHandler/exception"
	"github.com/iot-my-world/brain/pkg/security/role/setup/exception"
	humanUserAdministrator "github.com/iot-my-world/brain/pkg/user/human/administrator"
	humanUserRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	humanUserValidator "github.com/iot-my-world/brain/pkg/user/human/validator"
)

var CompanyAdmin = role.Role{
	Name:           "companyAdmin",
	APIPermissions: make([]apiPermission.Permission, 0),
	ViewPermissions: []viewPermission.Permission{
		viewPermission.PartyProfileEditing,

		viewPermission.PartyClient,
		viewPermission.PartyUser,

		viewPermission.LiveTrackingDashboard,
		viewPermission.HistoricalTrackingDashboard,

		viewPermission.TrackerSF001,
	},
}

var CompanyUser = role.Role{
	Name:           "companyUser",
	APIPermissions: make([]apiPermission.Permission, 0),
	ViewPermissions: []viewPermission.Permission{
		viewPermission.LiveTrackingDashboard,
		viewPermission.HistoricalTrackingDashboard,

		viewPermission.TrackerSF001,
	},
}

var ClientAdmin = role.Role{
	Name:           "clientAdmin",
	APIPermissions: make([]apiPermission.Permission, 0),
	ViewPermissions: []viewPermission.Permission{
		viewPermission.PartyProfileEditing,

		viewPermission.PartyUser,

		viewPermission.LiveTrackingDashboard,
		viewPermission.HistoricalTrackingDashboard,

		viewPermission.TrackerSF001,
	},
}

var ClientUser = role.Role{
	Name:           "clientUser",
	APIPermissions: make([]apiPermission.Permission, 0),
	ViewPermissions: []viewPermission.Permission{
		viewPermission.LiveTrackingDashboard,
		viewPermission.HistoricalTrackingDashboard,

		viewPermission.TrackerSF001,
	},
}

var initialRoles = func() []role.Role {

	rootAPIPermissions := make([]apiPermission.Permission, 0)

	// System RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, systemRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, systemRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, systemRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, systemRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, systemRecordHandler.ClientUserPermissions...)

	// Role RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, recordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, recordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, recordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, recordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, recordHandler.ClientUserPermissions...)

	// Permission Administrator
	rootAPIPermissions = append(rootAPIPermissions, administrator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, administrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, administrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, administrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, administrator.ClientUserPermissions...)

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

	// Sigbug Administrator
	rootAPIPermissions = append(rootAPIPermissions, sigbugAdministrator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, sigbugAdministrator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, sigbugAdministrator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, sigbugAdministrator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, sigbugAdministrator.ClientUserPermissions...)
	// Sigbug RecordHandler
	rootAPIPermissions = append(rootAPIPermissions, sigbugRecordHandler.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, sigbugRecordHandler.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, sigbugRecordHandler.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, sigbugRecordHandler.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, sigbugRecordHandler.ClientUserPermissions...)
	// Sigbug Validator
	rootAPIPermissions = append(rootAPIPermissions, sigbugValidator.SystemUserPermissions...)
	CompanyAdmin.APIPermissions = append(CompanyAdmin.APIPermissions, sigbugValidator.CompanyAdminUserPermissions...)
	CompanyUser.APIPermissions = append(CompanyUser.APIPermissions, sigbugValidator.CompanyUserPermissions...)
	ClientAdmin.APIPermissions = append(ClientAdmin.APIPermissions, sigbugValidator.ClientAdminUserPermissions...)
	ClientUser.APIPermissions = append(ClientUser.APIPermissions, sigbugValidator.ClientUserPermissions...)

	// Register roles here
	allRoles := []role.Role{
		ClientAdmin,
		ClientUser,
		CompanyAdmin,
		CompanyUser,
	}

	// The view permissions that root has
	rootViewPermissions := []viewPermission.Permission{
		viewPermission.PartyCompany,
		viewPermission.PartyClient,
		viewPermission.PartyUser,
		viewPermission.PartyAPIUser,

		viewPermission.LiveTrackingDashboard,
		viewPermission.HistoricalTrackingDashboard,

		viewPermission.TrackerSF001,
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

func InitialSetup(handler recordHandler.RecordHandler) error {
	for _, roleToCreate := range initialRoles {
		//Try and retrieve the record
		retrieveRoleResponse, err := handler.Retrieve(&recordHandler.RetrieveRequest{
			Identifier: name.Identifier{Name: roleToCreate.Name},
		})
		switch err.(type) {
		case roleRecordHandlerException.NotFound:
			// if role record does not exist yet, try and create it
			_, err := handler.Create(&recordHandler.CreateRequest{
				Role: roleToCreate,
			})
			if err != nil {
				return exception.InitialSetup{Reasons: []string{"creation error", err.Error()}}
			}
			log.Info("Initial Role Setup: Created Role: " + roleToCreate.Name)

		case nil:
			//Record Retrieved Successfully

			// Update Role Permissions If Necessary
			if !(roleToCreate.CompareAPIPermissions(retrieveRoleResponse.Role.APIPermissions) &&
				roleToCreate.CompareViewPermissions(retrieveRoleResponse.Role.ViewPermissions)) {
				// role permissions differ, try update role
				log.Info("Initial Role Setup: Role: " + roleToCreate.Name + " already exists. Updating Role API permissions.")
				if _, err := handler.Update(&recordHandler.UpdateRequest{
					Role:       roleToCreate,
					Identifier: id.Identifier{Id: retrieveRoleResponse.Role.Id},
				}); err != nil {
					return exception.InitialSetup{Reasons: []string{"update error", err.Error()}}
				}
			}

		default:
			// otherwise there was some retrieval error
			return exception.InitialSetup{Reasons: []string{"retrieval error", err.Error()}}
		}
	}

	return nil
}
