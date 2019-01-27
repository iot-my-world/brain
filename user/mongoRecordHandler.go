package user

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"strings"
)

type mongoRecordHandler struct{
	mongoSession *mgo.Session
	database, collection string
}

func NewMongoRecordHandler(mongoSession *mgo.Session, database, collection string) *mongoRecordHandler {

	setupUserRecords(mongoSession, database, collection)

	newUserMongoRecordHandler := mongoRecordHandler{
		mongoSession,
		database,
		collection,
	}

	//if err := initialUserSetup(&newUserMongoRecordHandler); err != nil {
	//	log.Fatal("Unable to complete initial user setup!", err)
	//}

	return &newUserMongoRecordHandler
}

func setupUserRecords(mongoSession *mgo.Session, database, collection string){
	//Initialise User collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	userCollection := mgoSesh.DB(database).C(collection)

	//Ensure Username Uniqueness
	username_unique := mgo.Index{
		Key: []string{"username"},
		Unique: true,
	}
	if err := userCollection.EnsureIndex(username_unique); err != nil {
		log.Fatal("Could not ensure uniqueness: ", err)
	}
}

func validateServiceReqData(r interface{}) error {
	var reasonsInvalid []string
	switch v := r.(type){
	case *CreateRequest:
		if v.NewUser.Name == "" {
			reasonsInvalid = append(reasonsInvalid, "Username cannot be blank")
		}
		if v.NewUser.Password == "" {
			reasonsInvalid = append(reasonsInvalid, "Password cannot be blank")
		}
	case *DeleteRequest:
		if v.Username == "" {
			reasonsInvalid = append(reasonsInvalid, "Username cannot be blank")
		}
	case *UpdateRequest:
		if v.UpdatedUser.Name == "" {
			reasonsInvalid = append(reasonsInvalid, "Username cannot be blank")
		}

	case *RetrieveAllRequest:
	default:
		log.Warn("NO CHECK CASE FOR THIS REQUEST!")
	}

	if len(reasonsInvalid) > 0 {
		return errors.New("Invalid Create Request: " + strings.Join(reasonsInvalid, ","))
	}
	return nil
}

func (u * mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {
	return nil
}

func (u *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {

	if err := validateServiceReqData(request); err != nil {
		return err
	}
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
		Name: request.NewUser.Name,
		Surname: request.NewUser.Surname,
		IDNo: request.NewUser.IDNo,

		// System Details
		Username: request.NewUser.Username,
		Password: pwdHash,
		SystemRole: request.NewUser.SystemRole,
		TagID: request.NewUser.TagID,

		// Business Details
		// TODO: Split out into "BusinessDetails" Struct
		BusinessRole: request.NewUser.BusinessRole,
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
	if err := validateServiceReqData(request); err != nil {
		return err
	}
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
	if err := validateServiceReqData(request); err != nil {
		return err
	}

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(u.database).C(u.collection)

	// Retrieve User
	retrievedUser := &User{}
	if err := userCollection.Find(bson.M{"username":request.UpdatedUser.Username}).One(retrievedUser); err != nil {
		return err
	}

	// Update fields:
	retrievedUser.TagID = request.UpdatedUser.TagID
	retrievedUser.Name = request.UpdatedUser.Name
	retrievedUser.Surname = request.UpdatedUser.Surname
	retrievedUser.IDNo = request.UpdatedUser.IDNo
	retrievedUser.Username = request.UpdatedUser.Username
	retrievedUser.SystemRole = request.UpdatedUser.SystemRole
	retrievedUser.BusinessRole = request.UpdatedUser.BusinessRole

	if err := userCollection.Update(bson.M{"username":request.UpdatedUser.Username}, retrievedUser); err != nil {
		log.Error("Unable to update user!", err)
		return err
	}

	response.User = *retrievedUser

	return nil
}

func (u *mongoRecordHandler) Delete(request *DeleteRequest, response *DeleteResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	userCollection := mgoSession.DB(u.database).C(u.collection)

	if err := userCollection.Remove(request); err != nil {
		return err
	}

	return nil
}

//func initialUserSetup(handler *mongoRecordHandler) error {
//	for _, userCreateRequest := range initialUsers {
//		//Try and retrieve the record
//		retrieveUserResponse := RetrieveResponse{}
//
//		if err := handler.Retrieve(&RetrieveRequest{Username:userCreateRequest.Username}, &retrieveUserResponse); err != nil {
//			//Unable to retrieve user record
//			//Try create user record
//			userCreateResponse := CreateResponse{}
//
//			if err := handler.Create(&userCreateRequest, &userCreateResponse); err != nil {
//				log.Fatal("Unable to Complete Initial User Setup!", "Could Not Create User: " + userCreateRequest.Username)
//			}
//			if !userCreateResponse.Success {
//				log.Fatal("Unable to Complete Initial User Setup!", "Could Not Create User: " + userCreateRequest.Username, userCreateResponse.Reasons)
//			}
//			log.Info("Initial User Setup: Created User: " + userCreateRequest.Username)
//			continue
//		}
//
//		//User Record Retrieved Successfully, update user record
//		log.Info("Initial User Setup: User " + userCreateRequest.Username + " already exists. Updating User.")
//		updateUserResponse := UpdateResponse{}
//		if err := handler.Update(&UpdateRequest{userCreateRequest.Username, userCreateRequest.Password, userCreateRequest.SystemRole, userCreateRequest.BusinessRole}, &updateUserResponse); err != nil {
//			log.Fatal("Unable to Complete Initial User Setup! Error updating User record!")
//		}
//		if !updateUserResponse.Success {
//			log.Fatal("Unable to Complete Initial User Setup!", updateUserResponse.Reasons)
//		}
//	}
//
//	return nil
//}

//func (u *mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {
//
//	response.Reasons, response.Success  = validateServiceReqData(request)
//	if !response.Success {
//		return nil
//	}
//
//	mgoSession := u.mongoSession.Copy()
//	defer mgoSession.Close()
//
//	userCollection := mgoSession.DB(u.database).C(u.collection)
//	retrievedUser := User{}
//	err := userCollection.Find(bson.M{"username":request.Username}).One(&retrievedUser)
//	if err != nil {
//		//log.Error("Unable to retrieve user", err)
//		return err
//	}
//	response.User = retrievedUser
//	response.Success = true
//	return nil
//}
//
//func (u *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {
//
//	response.Reasons, response.Success  = validateServiceReqData(request)
//	if !response.Success {
//		return nil
//	}
//	mgoSession := u.mongoSession.Copy()
//	defer mgoSession.Close()
//
//	userCollection := mgoSession.DB(u.database).C(u.collection)
//
//	pwdHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
//	if err != nil {
//		log.Error("Error While hashing Password!", err)
//		return err
//	}
//
//	//TODO: Confirm that the role that we are attempting to update the user with is a real role
//
//	err = userCollection.Update(
//		bson.M{"username":request.Username},
//		&User{
//			Username: request.Username,
//			Password: pwdHash,
//			SystemRole: request.SystemRole,
//			BusinessRole: request.BusinessRole,
//		},
//	)
//
//	if err != nil {
//		log.Error("Could not update user! ", err)
//		return err //TODO: Translate Unknown error
//	}
//	response.Success = true
//	return nil
//}
//
//func (u *mongoRecordHandler) Delete(request *DeleteRequest, response *DeleteResponse) error {
//	return nil
//}
