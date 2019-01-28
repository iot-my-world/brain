package user

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/iotTracker/brain/party"
	"gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/search/identifiers/name"
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

func (mrh *mongoRecordHandler) ValidateRetrieveRequest(request *RetrieveRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Identifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return exception.RequestInvalid{Reasons: reasonsInvalid}
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
		return err
	}

	response.User = userRecord
	return nil
}

func (mrh *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error While hashing Password!", err)
		return err
	}

	userToInsert := &party.User{
		// Personal Details
		// TODO: Split out into "PersonalDetails" Struct
		Name:    request.NewUser.Name,
		Surname: request.NewUser.Surname,
		IDNo:    request.NewUser.IDNo,

		// System Details
		Username:   request.NewUser.Username,
		Password:   pwdHash,
		SystemRole: request.NewUser.SystemRole,
	}

	err = userCollection.Insert(userToInsert)
	if err != nil {
		log.Error("Could not create user! ", err)
		return err //TODO: Translate Unknown error
	}
	response.User = *userToInsert

	return nil
}

func (mrh *mongoRecordHandler) RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error {
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	var records []party.User

	if err := userCollection.Find(bson.M{}).All(&records); err != nil {
		return err
	}

	response.UserRecords = records
	return nil
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

func (mrh *mongoRecordHandler) Delete(request *DeleteRequest, response *DeleteResponse) error {
	mgoSession := mrh.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(mrh.database).C(mrh.collection)

	if err := userCollection.Remove(request); err != nil {
		return err
	}

	return nil
}

func initialUserSetup(handler *mongoRecordHandler) error {
	for _, newUser := range initialUsers {
		//Try and retrieve the new user record
		retrieveUserResponse := RetrieveResponse{}

		if err := handler.Retrieve(&RetrieveRequest{
			Identifier: name.Identifier(newUser.Name),
		}, &retrieveUserResponse); err != nil {
			// user could not be found
			if err != mgo.ErrNotFound {
				log.Fatal("Unable to Complete Initial User Setup!", "Could Not Find User: "+newUser.Username+err.Error())
			}
			// User Record does not exist
			//Try create user record
			userCreateResponse := CreateResponse{}

			if err := handler.Create(&CreateRequest{
				NewUser: newUser,
			}, &userCreateResponse); err != nil {
				log.Fatal("Unable to Complete Initial User Setup!", "Could Not Create User: "+newUser.Username)
			}
			if err != nil {
				log.Fatal("Unable to Complete Initial User Setup!", "Could Not Create User: "+newUser.Username, err.Error())
			}
			log.Info("Initial User Setup: Created User: " + newUser.Username)
			continue
		}

		//User Record Retrieved Successfully, update user record
		log.Info("Initial User Setup: User " + newUser.Username + " already exists. Updating User.")
	}

	return nil
}
