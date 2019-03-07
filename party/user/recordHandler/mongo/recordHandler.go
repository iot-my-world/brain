package mongo

import (
	"fmt"
	"github.com/satori/go.uuid"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/party/user"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/party/user/recordHandler/exception"
	userSetup "gitlab.com/iotTracker/brain/party/user/setup"
	partyRegistrar "gitlab.com/iotTracker/brain/party/registrar"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/api"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/username"
)

type mongoRecordHandler struct {
	mongoSession   *mgo.Session
	database       string
	collection     string
	ignoredReasons map[api.Method]reasonInvalid.IgnoredReasonsInvalid
}

func New(
	mongoSession *mgo.Session,
	database,
	collection string,
) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	ignoredReasons := map[api.Method]reasonInvalid.IgnoredReasonsInvalid{
		userRecordHandler.Create: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"password": {
					reasonInvalid.Blank,
				},
			},
		},

		partyRegistrar.InviteCompanyAdminUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"name": {
					reasonInvalid.Blank,
				},
				"surname": {
					reasonInvalid.Blank,
				},
				"username": {
					reasonInvalid.Blank,
				},
				"password": {
					reasonInvalid.Blank,
				},
			},
		},

		partyRegistrar.RegisterCompanyAdminUser: {
			ReasonsInvalid: map[string][]reasonInvalid.Type{
				"id": {
					reasonInvalid.Blank,
				},
				"password": {
					reasonInvalid.Blank,
				},
			},
		},
	}

	newUserMongoRecordHandler := mongoRecordHandler{
		mongoSession:   mongoSession,
		database:       database,
		collection:     collection,
		ignoredReasons: ignoredReasons,
	}

	if err := userSetup.InitialSetup(&newUserMongoRecordHandler); err != nil {
		log.Fatal("Unable to complete initial user setup!", err.Error())
	}

	return &newUserMongoRecordHandler
}

func setupIndices(mongoSession *mgo.Session, database, collection string) {
	//Initialise User collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	userCollection := mgoSesh.DB(database).C(collection)

	// Ensure id uniqueness
	idUnique := mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	}
	if err := userCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure id uniqueness: ", err)
	}

	//Ensure username Uniqueness
	usernameUnique := mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}
	if err := userCollection.EnsureIndex(usernameUnique); err != nil {
		log.Fatal("Could not ensure username uniqueness: ", err)
	}

	//Ensure emailAddress Uniqueness
	emailAddressUnique := mgo.Index{
		Key:    []string{"emailAddress"},
		Unique: true,
	}
	if err := userCollection.EnsureIndex(emailAddressUnique); err != nil {
		log.Fatal("Could not ensure email address uniqueness: ", err)
	}
}

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *userRecordHandler.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		if request.Claims.PartyDetails().PartyType != party.System {
			// If the user validating for a create is not root then the user's party must be the the party of the user
			if request.User.PartyType != request.Claims.PartyDetails().PartyType {
				reasonsInvalid = append(reasonsInvalid, "partyType must be submitting party's type")
			}
			if request.User.PartyId.Id != request.Claims.PartyDetails().PartyId.Id {
				reasonsInvalid = append(reasonsInvalid, "partyId must be submitting party's id")
			}
		}

		// Validate the new user
		userValidateResponse := userRecordHandler.ValidateResponse{}
		err := mrh.Validate(&userRecordHandler.ValidateRequest{
			Claims: request.Claims,
			User:   request.User,
			Method: userRecordHandler.Create,
		}, &userValidateResponse)
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "unable to validate newUser")
		} else {
			for _, reason := range userValidateResponse.ReasonsInvalid {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Create(request *userRecordHandler.CreateRequest, response *userRecordHandler.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	newId, err := uuid.NewV4()
	if err != nil {
		return brainException.UUIDGeneration{Reasons: []string{err.Error()}}
	}
	request.User.Id = newId.String()

	if err := userCollection.Insert(request.User); err != nil {
		return userRecordHandlerException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.User = request.User
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *userRecordHandler.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *userRecordHandler.RetrieveRequest, response *userRecordHandler.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var userRecord user.User

	filter := user.ContextualiseFilter(request.Identifier.ToFilter(), request.Claims)
	if err := userCollection.Find(filter).One(&userRecord); err != nil {
		if err == mgo.ErrNotFound {
			return userRecordHandlerException.NotFound{}
		} else {
			return brainException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.User = userRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *userRecordHandler.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *userRecordHandler.UpdateRequest, response *userRecordHandler.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve User
	retrieveUserResponse := userRecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&userRecordHandler.RetrieveRequest{Identifier: request.Identifier}, &retrieveUserResponse); err != nil {
		return userRecordHandlerException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveUserResponse.User.Id = request.User.Id // cannot update ever
	retrieveUserResponse.User.Name = request.User.Name
	retrieveUserResponse.User.Surname = request.User.Surname
	// retrieveUserResponse.User.Username = request.User.Username // cannot update yet
	// retrieveUserResponse.User.EmailAddress = request.User.EmailAddress // cannot update yet
	retrieveUserResponse.User.Password = request.User.Password
	retrieveUserResponse.User.Roles = request.User.Roles
	// retrieveUserResponse.User.PartyType = request.User.PartyType // cannot update yet
	// retrieveUserResponse.User.PartyId = request.User.PartyId // cannot update yet

	if err := userCollection.Update(request.Identifier.ToFilter(), retrieveUserResponse.User); err != nil {
		return userRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.User = retrieveUserResponse.User

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *userRecordHandler.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *userRecordHandler.DeleteRequest, response *userRecordHandler.DeleteResponse) error {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := userCollection.Remove(request.Identifier.ToFilter()); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *userRecordHandler.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Validate(request *userRecordHandler.ValidateRequest, response *userRecordHandler.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	userToValidate := &request.User

	if (*userToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*userToValidate).Id,
		})
	}

	if (*userToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Name,
		})
	}

	if (*userToValidate).Surname == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "surname",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Name,
		})
	}

	if (*userToValidate).Username == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "username",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Username,
		})
	}

	if (*userToValidate).EmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "emailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).EmailAddress,
		})
	}

	if len((*userToValidate).Password) == 0 {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "password",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Password,
		})
	}

	if (*userToValidate).ParentPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).ParentPartyType,
		})
	}

	if (*userToValidate).ParentId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).PartyId,
		})
	}

	if (*userToValidate).PartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "partyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).PartyType,
		})
	}

	if (*userToValidate).PartyId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "partyId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).PartyId,
		})
	}

	switch request.Method {
	case userRecordHandler.Create, partyRegistrar.RegisterCompanyAdminUser:
		// Check if the users username has already been assigned to another user
		if (*userToValidate).Username != "" {
			if err := mrh.Retrieve(&userRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: username.Identifier{
					Username: (*userToValidate).Username,
				},
			},
				&userRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case userRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "username",
						Type:  reasonInvalid.Unknown,
						Help:  "retrieve failed",
						Data:  (*userToValidate).Username,
					})
				}
			} else {
				// there was no error, this email address is already taken by another user
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "username",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*userToValidate).Username,
				})
			}
		}
		fallthrough

	case partyRegistrar.InviteCompanyAdminUser:
		// Check if the users email has already been assigned to another user
		if (*userToValidate).EmailAddress != "" {
			if err := mrh.Retrieve(&userRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: emailAddress.Identifier{
					EmailAddress: (*userToValidate).EmailAddress,
				},
			},
				&userRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case userRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "emailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "retrieve failed",
						Data:  (*userToValidate).EmailAddress,
					})
				}
			} else {
				// there was no error, this email address is already taken by another user
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "emailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*userToValidate).EmailAddress,
				})
			}
		}
	}

	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Ignore reasons applicable to method if relevant
	if mrh.ignoredReasons[request.Method].ReasonsInvalid != nil {
		for _, reason := range allReasonsInvalid {
			if !mrh.ignoredReasons[request.Method].CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}
	}

	response.ReasonsInvalid = returnedReasonsInvalid
	return nil
}

func (mrh *mongoRecordHandler) ValidateChangePasswordRequest(request *userRecordHandler.ChangePasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) ChangePassword(request *userRecordHandler.ChangePasswordRequest, response *userRecordHandler.ChangePasswordResponse) error {
	if err := mrh.ValidateChangePasswordRequest(request); err != nil {
		return err
	}

	// Retrieve User
	retrieveUserResponse := userRecordHandler.RetrieveResponse{}
	if err := mrh.Retrieve(&userRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.Identifier,
	}, &retrieveUserResponse); err != nil {
		return userRecordHandlerException.ChangePassword{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return userRecordHandlerException.ChangePassword{Reasons: []string{"hashing password", err.Error()}}
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// update user
	retrieveUserResponse.User.Password = pwdHash

	if err := userCollection.Update(request.Identifier.ToFilter(), retrieveUserResponse.User); err != nil {
		return userRecordHandlerException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.User = retrieveUserResponse.User

	return nil
}
