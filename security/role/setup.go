package role

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/identifier/name"
	"gitlab.com/iotTracker/brain/security/permission"
	roleException "gitlab.com/iotTracker/brain/security/role/exception"
)

var initialRoles = func() []Role {

	// Register roles here
	allRoles := []Role{
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
	}

	// Create root role and apply permissions of all other roles to root
	for _, role := range allRoles {
		rootPermissions = append(rootPermissions, role.Permissions...)
	}
	root := Role{
		Name:        "root",
		Permissions: rootPermissions,
	}
	return append([]Role{root}, allRoles...)
}()

// Create Roles here

var owner = Role{
	Name: "client",
	Permissions: []permission.Permission{
		"User.Retrieve",
	},
}

var admin = Role{
	Name:        "admin",
	Permissions: []permission.Permission{},
}

func InitialSetup(handler RecordHandler) error {
	for _, roleToCreate := range initialRoles {
		//Try and retrieve the record
		retrieveRoleResponse := RetrieveResponse{}
		err := handler.Retrieve(&RetrieveRequest{Identifier: name.Identifier{Name: roleToCreate.Name}}, &retrieveRoleResponse)

		switch err.(type) {
		case roleException.NotFound:
			// if role record does not exist yet, try and create it
			createRoleResponse := CreateResponse{}
			if err := handler.Create(&CreateRequest{Role: roleToCreate}, &createRoleResponse); err != nil {
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
				if err := handler.Update(&UpdateRequest{Role: roleToCreate}, &UpdateResponse{}); err != nil {
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
