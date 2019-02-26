package setup

import (
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandlerException "gitlab.com/iotTracker/brain/party/user/recordHandler/exception"
	userSetupException "gitlab.com/iotTracker/brain/party/user/setup/exception"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/identifier/username"
	"os"
	"io/ioutil"
)

type newUser struct {
	user     user.User
	password string
}

//var initialUsers = []newUser{
//	{
//		user: user.User{
//			Name:    "root",
//			Surname: "root",
//
//			Username:     "root",
//			EmailAddress: "root@root.com",
//			// Password: set later with hashing
//			Roles: []string{"root"},
//
//			PartyType: party.System,
//			PartyId:   id.Identifier{Id: "root"},
//		},
//		password: "12345",
//	},
//}

var initialUsers = make([]newUser, 0)

func consumePasswordFile(location string) ([]byte, error) {
	if _, err := os.Stat(location); err != nil {
		return nil, err
	}
	// read the file
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	// remove the file
	if err := os.Remove(location); err != nil {
		return nil, err
	}
	// return the data
	return data, nil
}

func InitialSetup(handler userRecordHandler.RecordHandler) error {
	for _, newUser := range initialUsers {
		//Try and retrieve the new user record
		retrieveUserResponse := userRecordHandler.RetrieveResponse{}
		err := handler.Retrieve(&userRecordHandler.RetrieveRequest{Identifier: username.Identifier{Username: newUser.user.Username}}, &retrieveUserResponse)

		switch err.(type) {
		case userRecordHandlerException.NotFound:
			// if user record does not exist yet, try and create it
			userCreateResponse := userRecordHandler.CreateResponse{}
			if err := handler.Create(&userRecordHandler.CreateRequest{User: newUser.user}, &userCreateResponse); err != nil {
				return userSetupException.InitialSetup{Reasons: []string{"creation error", err.Error()}}
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
				return userSetupException.InitialSetup{Reasons: []string{"update error", err.Error()}}
			}

		default:
			// otherwise there was some retrieval error
			return userSetupException.InitialSetup{Reasons: []string{"retrieval error", err.Error()}}
		}

		// creation or update done, update password
		// creation done, change password
		userChangePasswordResponse := userRecordHandler.ChangePasswordResponse{}
		if err := handler.ChangePassword(&userRecordHandler.ChangePasswordRequest{
			Identifier:  username.Identifier{Username: newUser.user.Username},
			NewPassword: newUser.password,
		}, &userChangePasswordResponse); err != nil {
			return userSetupException.InitialSetup{Reasons: []string{"change password error", err.Error()}}
		}
	}

	return nil
}
