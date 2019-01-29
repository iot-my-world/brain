package user

import (
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifiers/name"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	"gitlab.com/iotTracker/brain/log"
)

var initialUsers = []party.NewUser{
	{
		Name:     "root",
		Username: "root",
		Password: "123",
	},
}

func initialUserSetup(handler *mongoRecordHandler) error {
	for _, newUser := range initialUsers {
		//Try and retrieve the new user record
		retrieveUserResponse := RetrieveResponse{}
		err := handler.Retrieve(&RetrieveRequest{Identifier: name.Identifier(newUser.Name)}, &retrieveUserResponse)

		switch err.(type) {
		case userException.NotFound:
			// if user record does not exist yet, try and create it
			userCreateResponse := CreateResponse{}
			if err := handler.Create(&CreateRequest{NewUser: newUser}, &userCreateResponse); err != nil {
				return userException.InitialSetup{Reasons: []string{"creation failure", err.Error()}}
			}
			log.Info("Initial User Setup: Created User: " + newUser.Username)
			continue // user created successfully, continue

		case nil:
			// no error, user was retrieved successfully
			log.Info("Initial User Setup: User " + newUser.Username + " already exists. Updating User.")
			userUpdateResponse := UpdateResponse{}
			err := handler.Update(&UpdateRequest{}, &userUpdateResponse)

		default:
			// otherwise there was some retrieval error
			return userException.InitialSetup{Reasons: []string{"retrieval failure", err.Error()}}
		}
	}

	return nil
}
