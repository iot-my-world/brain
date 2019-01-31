package role

import (
	"gitlab.com/iotTracker/brain/security"
	"gitlab.com/iotTracker/brain/search/identifiers/name"
	roleException "gitlab.com/iotTracker/brain/security/role/exception"
	"gitlab.com/iotTracker/brain/log"
)

var initialRoles = func() []security.Role {

	// Register roles here
	allRoles := []security.Role{
		owner,
		admin,
	}

	//Register additional root permissions here
	rootPermissions := []security.Permission{
		"Role.Create",
		"Role.Retrieve",
		"Role.Update",
		"Role.Delete",
	}

	// Create root role and apply permissions of all other roles to root
	for _, role := range allRoles {
		rootPermissions = append(rootPermissions, role.Permissions...)
	}
	root := security.Role{
		Name:        "root",
		Permissions: rootPermissions,
	}
	return append([]security.Role{root}, allRoles...)
}()

// Create Roles here

var owner = security.Role{
	Name: "client",
	Permissions: []security.Permission{
		"User.Retrieve",
	},
}

var admin = security.Role{
	Name: "admin",
	Permissions: []security.Permission{
	},
}

func initialRoleSetup(handler *mongoRecordHandler) error {
	for _, roleToCreate := range initialRoles {
		//Try and retrieve the record
		retrieveRoleResponse := RetrieveResponse{}
		err := handler.Retrieve(&RetrieveRequest{Identifier: name.Identifier(roleToCreate.Name)}, &retrieveRoleResponse)

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
