package user

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/iotTracker/brain/party"
	globalException "gitlab.com/iotTracker/brain/exception"
	userException "gitlab.com/iotTracker/brain/party/user/exception"
	"fmt"
	"gitlab.com/iotTracker/brain/validate"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type mongoRecordHandler struct {
	mongoSession         *mgo.Session
	database, collection string
}

func NewMongoRecordHandler(mongoSession *mgo.Session, database, collection string) *mongoRecordHandler {

	setupUserRecords(mongoSession, database, collection)

	newUserMongoRecordHandler := mongoRecordHandler{
		mongoSession,
		database,
		collection,
	}

	if err := initialUserSetup(&newUserMongoRecordHandler); err != nil {
		log.Fatal("Unable to complete initial user setup!", err)
	}

	return &newUserMongoRecordHandler
}

func setupUserRecords(mongoSession *mgo.Session, database, collection string) {
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

func (mrh *mongoRecordHandler) ValidateCreateRequest(request *CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {
	if err := mrh.ValidateCreateRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return userException.Create{Reasons: []string{"hashing password", err.Error()}}
	}

	userToInsert := &party.User{
		// Personal Details
		Name:    request.NewUser.Name,
		Surname: request.NewUser.Surname,

		// System Details
		Username:     request.NewUser.Username,
		EmailAddress: request.NewUser.Password,
		Password:     pwdHash,
		SystemRole:   request.NewUser.SystemRole,
	}

	err = userCollection.Insert(userToInsert)

	if err != nil {
		return userException.Create{Reasons: []string{"inserting record", err.Error()}}
	}

	response.User = *userToInsert
	return nil
}

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {
	if err := mrh.ValidateRetrieveRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var userRecord party.User

	if err := userCollection.Find(request.Identifier.ToMap()).One(&userRecord); err != nil {
		if err == mgo.ErrNotFound {
			return userException.NotFound{}
		} else {
			return globalException.Unexpected{Reasons: []string{err.Error()}}
		}
	}

	response.User = userRecord
	return nil
}

func (mrh *mongoRecordHandler) ValidateUpdateRequest(request *CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	// Retrieve User
	retrievedUser := &party.User{}
	if err := userCollection.Find(bson.M{"username": request.UpdatedUser.Username}).One(retrievedUser); err != nil {
		return err
	}

	// Update fields:
	retrievedUser.Name = request.UpdatedUser.Name
	retrievedUser.Surname = request.UpdatedUser.Surname
	retrievedUser.IDNo = request.UpdatedUser.IDNo
	retrievedUser.Username = request.UpdatedUser.Username
	retrievedUser.SystemRole = request.UpdatedUser.SystemRole

	if err := userCollection.Update(bson.M{"username": request.UpdatedUser.Username}, retrievedUser); err != nil {
		log.Error("Unable to update user!", err)
		return err
	}

	response.User = *retrievedUser

	return nil
}

func (mrh *mongoRecordHandler) ValidateDeleteRequest(request *DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !IsValidIdentifier(request.Identifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.Identifier.Type()))
		}
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Delete(request *DeleteRequest, response *DeleteResponse) error {
	if err := mrh.ValidateDeleteRequest(request); err != nil {
		return err
	}

	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := userCollection.Remove(request.Identifier.ToMap()); err != nil {
		return err
	}

	return nil
}

func (mrh *mongoRecordHandler) ValidateValidateRequest(request *ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (mrh *mongoRecordHandler) Validate(request *ValidateRequest, response *ValidateResponse) error {
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

	if (*userToValidate).Name == "" {
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
