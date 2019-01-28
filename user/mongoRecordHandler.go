package user

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
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

	//Ensure Username Uniqueness
	username_unique := mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}
	if err := userCollection.EnsureIndex(username_unique); err != nil {
		log.Fatal("Could not ensure uniqueness: ", err)
	}
}

func (u *mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {
	return nil
}

func (u *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(u.database).C(u.collection)

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.NewUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error While hashing Password!", err)
		return err
	}

	userToInsert := &User{
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

func (u *mongoRecordHandler) RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error {
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(u.database).C(u.collection)

	var records []User

	if err := userCollection.Find(bson.M{}).All(&records); err != nil {
		return err
	}

	response.UserRecords = records
	return nil
}

func (u *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(u.database).C(u.collection)

	// Retrieve User
	retrievedUser := &User{}
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

func (u *mongoRecordHandler) Delete(request *DeleteRequest, response *DeleteResponse) error {
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(u.database).C(u.collection)

	if err := userCollection.Remove(request); err != nil {
		return err
	}

	return nil
}

func initialUserSetup(handler *mongoRecordHandler) error {
	for _, newUser := range initialUsers {
		//Try and retrieve the new user record
		retrieveUserResponse := RetrieveResponse{}

		if err := handler.Retrieve(&RetrieveRequest{Username: newUser.Username}, &retrieveUserResponse); err != nil {
			//Unable to retrieve user record
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
