package setup

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/security/permission"
	"gitlab.com/iotTracker/brain/security/role"
	roleException "gitlab.com/iotTracker/brain/security/role/exception"
	roleRecordHandler "gitlab.com/iotTracker/brain/security/role/recordHandler"
)

var initialRoles = func() []role.Role {

	// Register roles here
	allRoles := []role.Role{
		owner,
		admin,
	}

	//Register additional root permissions here
	rootPermissions := []permission.Permission{
		permission.RoleCreate,
		permission.RoleRetrieve,
		permission.RoleUpdate,
		permission.RoleDelete,
		permission.CompanyRecordHandlerCreate,
		permission.CompanyRecordHandlerRetrieve,
		permission.CompanyRecordHandlerUpdate,
		permission.CompanyRecordHandlerDelete,
		permission.CompanyRecordHandlerValidate,
		permission.CompanyRecordHandlerCollect,
		permission.PartyRegistrarInviteCompanyAdminUser,
	}

	// Create root role and apply permissions of all other roles to root
	for _, role := range allRoles {
		rootPermissions = append(rootPermissions, role.Permissions...)
	}
	root := role.Role{
		Name:        "root",
		Permissions: rootPermissions,
	}
	return append([]role.Role{root}, allRoles...)
}()

// Create Roles here

var owner = role.Role{
	Name: "client",
	Permissions: []permission.Permission{
		"User.Retrieve",
	},
}

var admin = role.Role{
	Name:        "admin",
	Permissions: []permission.Permission{},
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
			// no error, role was retrieved successfully
			//Record Retrieved Successfully
			if roleToCreate.ComparePermissions(retrieveRoleResponse.Role.Permissions) {
				// no difference in role permissions, do nothing
				log.Info("Initial Role Setup: Role " + retrieveRoleResponse.Role.Name + " already exists and permissions correct.")
			} else {
				// role permissions differ, try update role
				log.Info("Initial Role Setup: Role: " + roleToCreate.Name + " already exists. Updating Role permissions.")
				if err := handler.Update(&roleRecordHandler.UpdateRequest{Role: roleToCreate}, &roleRecordHandler.UpdateResponse{}); err != nil {
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
