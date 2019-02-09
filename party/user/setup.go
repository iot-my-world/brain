package user

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
)

type newUser struct {
	user     User
	password string
}

var initialUsers = []newUser{
	{
		user: User{
			Name:    "root",
			Surname: "root",

			Username:     "root",
			EmailAddress: "root@root.com",
			// Password: set later with hashing
			Roles: []string{"root"},

			PartyType: party.System,
			// PartyId: no associated party for root
		},
		password: "12345",
	},
}

func InitialSetup(handler RecordHandler) error {
	for _, newUser := range initialUsers {
		//Try and retrieve the new user record
		retrieveUserResponse := RetrieveResponse{}
		err := handler.Retrieve(&RetrieveRequest{Identifier: username.Identifier{Username: newUser.user.Username}}, &retrieveUserResponse)

		switch err.(type) {
		case userException.NotFound:
			// if user record does not exist yet, try and create it
			userCreateResponse := CreateResponse{}
			if err := handler.Create(&CreateRequest{User: newUser.user}, &userCreateResponse); err != nil {
				return userException.InitialSetup{Reasons: []string{"creation error", err.Error()}}
			}
			log.Info("Initial User Setup: Created User: " + newUser.user.Username)

		case nil:
			// no error, user was retrieved successfully
			log.Info("Initial User Setup: User " + newUser.user.Username + " already exists. Updating User.")
			userUpdateResponse := UpdateResponse{}
			if err := handler.Update(&UpdateRequest{
				Identifier: id.Identifier{Id: retrieveUserResponse.User.Id},
				User:       newUser.user,
			}, &userUpdateResponse); err != nil {
				return userException.InitialSetup{Reasons: []string{"update error", err.Error()}}
			}

		default:
			// otherwise there was some retrieval error
			return userException.InitialSetup{Reasons: []string{"retrieval error", err.Error()}}
		}

		// creation or update done, update password
		// creation done, change password
		userChangePasswordResponse := ChangePasswordResponse{}
		if err := handler.ChangePassword(&ChangePasswordRequest{
			Identifier:  username.Identifier{Username: newUser.user.Username},
			NewPassword: newUser.password,
		}, &userChangePasswordResponse); err != nil {
			return userException.InitialSetup{Reasons: []string{"change password error", err.Error()}}
		}
	}

	return nil
}
