package mongo

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"gopkg.in/mgo.v2/bson"
	globalException "gitlab.com/iotTracker/brain/exception"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	"fmt"
	"gitlab.com/iotTracker/brain/validate"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
	"golang.org/x/crypto/bcrypt"
	"gitlab.com/iotTracker/brain/party/user"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database             string
	collection           string
	createIgnoredReasons validate.IgnoredReasonsInvalid
}

func New(
	mongoSession *mgo.Session,
	database string,
	collection string,
) *mongoRecordHandler {

	setupIndices(mongoSession, database, collection)

	createIgnoredReasons := validate.IgnoredReasonsInvalid{
		ReasonsInvalid: map[string][]reasonInvalid.Type{
			"id": {
				reasonInvalid.Blank,
			},
		},
	}

	newUserMongoRecordHandler := mongoRecordHandler{
		mongoSession:         mongoSession,
		database:             database,
		collection:           collection,
		createIgnoredReasons: createIgnoredReasons,
	}

	if err := user.InitialSetup(&newUserMongoRecordHandler); err != nil {
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

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *user.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// Validate the new user
	userValidateResponse := user.ValidateResponse{}

	err := mrh.Validate(&user.ValidateRequest{User: request.User}, &userValidateResponse)
	if err != nil {
		reasonsInvalid = append(reasonsInvalid, "unable to validate newUser")
	} else {
		for _, reason := range userValidateResponse.ReasonsInvalid {
			if !mrh.createIgnoredReasons.CanIgnore(reason) {
				reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("%s - %s", reason.Field, reason.Type))
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Create(request *user.CreateRequest, response *user.CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	request.User.Id = bson.NewObjectId().Hex()

	if err := userCollection.Insert(request.User); err != nil {
		return userException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.User = request.User
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *user.RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *user.RetrieveRequest, response *user.RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var userRecord user.User

	if err := userCollection.Find(request.Identifier.ToFilter()).One(&userRecord); err != nil {
		if err == mgo.ErrNotFound {
			return userException.NotFound{}
		} else {
			return globalException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.User = userRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *user.UpdateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *user.UpdateRequest, response *user.UpdateResponse) error {
	if err := mrh.ValidateUpdateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve User
	retrieveUserResponse := user.RetrieveResponse{}
	if err := mrh.Retrieve(&user.RetrieveRequest{Identifier: request.Identifier}, &retrieveUserResponse); err != nil {
		return userException.Update{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Update fields:
	// retrieveUserResponse.User.Id = request.User.Id // cannot update ever
	retrieveUserResponse.User.Name = request.User.Name
	retrieveUserResponse.User.Surname = request.User.Surname
	// retrieveUserResponse.User.Username = request.User.Username // cannot update yet
	// retrieveUserResponse.User.EmailAddress = request.User.EmailAddress // cannot update yet
	retrieveUserResponse.User.Password = request.User.Password
	retrieveUserResponse.User.Roles = request.User.Roles

	if err := userCollection.Update(request.Identifier.ToFilter(), retrieveUserResponse.User); err != nil {
		return userException.Update{Reasons: []string{"updating record", err.Error()}}
	}

	response.User = retrieveUserResponse.User

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *user.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *user.DeleteRequest, response *user.DeleteResponse) error {
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

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *user.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Validate(request *user.ValidateRequest, response *user.ValidateResponse) error {
	if err := mrh.ValidateValidateRequest(request); err != nil {
		return err
	}

	reasonsInvalid := make([]validate.ReasonInvalid, 0)
	userToValidate := &request.User

	if (*userToValidate).Id == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*userToValidate).Id,
		})
	}

	if (*userToValidate).Name == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Name,
		})
	}

	if (*userToValidate).Surname == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "surname",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Name,
		})
	}

	if (*userToValidate).Username == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "username",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).Username,
		})
	}

	if (*userToValidate).EmailAddress == "" {
		reasonsInvalid = append(reasonsInvalid, validate.ReasonInvalid{
			Field: "emailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*userToValidate).EmailAddress,
		})
	}

	response.ReasonsInvalid = reasonsInvalid
	return nil
}

func (mrh *mongoRecordHandler) ValidateChangePasswordRequest(request *user.ChangePasswordRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) ChangePassword(request *user.ChangePasswordRequest, response *user.ChangePasswordResponse) error {
	if err := mrh.ValidateChangePasswordRequest(request); err != nil {
		return err
	}

	// Retrieve User
	retrieveUserResponse := user.RetrieveResponse{}
	if err := mrh.Retrieve(&user.RetrieveRequest{Identifier: request.Identifier}, &retrieveUserResponse); err != nil {
		return userException.ChangePassword{Reasons: []string{"retrieving record", err.Error()}}
	}

	// Hash the new Password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return userException.ChangePassword{Reasons: []string{"hashing password", err.Error()}}
	}

	// update user
	retrieveUserResponse.User.Password = pwdHash
	updateUserResponse := user.UpdateResponse{}
	if err := mrh.Update(&user.UpdateRequest{Identifier: request.Identifier, User: retrieveUserResponse.User}, &updateUserResponse);
		err != nil {
		return userException.ChangePassword{Reasons: []string{"updating user", err.Error()}}
	}

	response.User = retrieveUserResponse.User

	return nil
}
