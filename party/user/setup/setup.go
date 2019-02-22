package setup

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/party/user"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
)

type newUser struct {
	user     user.User
	password string
}

var initialUsers = []newUser{
	{
		user: user.User{
			Name:    "root",
			Surname: "root",

			Username:     "root",
			EmailAddress: "root@root.com",
			// Password: set later with hashing
			Roles: []string{"root"},

			PartyType: party.System,
			PartyId:   id.Identifier{Id: "root"},
		},
		password: "12345",
	},
}

func InitialSetup(handler userRecordHandler.RecordHandler) error {
	for _, newUser := range initialUsers {
		//Try and retrieve the new user record
		retrieveUserResponse := userRecordHandler.RetrieveResponse{}
		err := handler.Retrieve(&userRecordHandler.RetrieveRequest{Identifier: username.Identifier{Username: newUser.user.Username}}, &retrieveUserResponse)

		switch err.(type) {
		case userException.NotFound:
			// if user record does not exist yet, try and create it
			userCreateResponse := userRecordHandler.CreateResponse{}
			if err := handler.Create(&userRecordHandler.CreateRequest{User: newUser.user}, &userCreateResponse); err != nil {
				return userException.InitialSetup{Reasons: []string{"creation error", err.Error()}}
			}
			log.Info("Initial User Setup: Created User: " + newUser.user.Username)

		case nil:
			// no error, user was retrieved successfully
			log.Info("Initial User Setup: User " + newUser.user.Username + " already exists. Updating User.")
			userUpdateResponse := userRecordHandler.UpdateResponse{}
			if err := handler.Update(&userRecordHandler.UpdateRequest{
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
		userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
		if err := handler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
			Identifier:  username.Identifier{Username: newUser.user.Username},
			NewPassword: newUser.password,
		}, &userChangePasswordResponse); err != nil {
			return userException.InitialSetup{Reasons: []string{"change password error", err.Error()}}
		}
	}

	return nil
}
