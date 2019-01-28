package systemRole

import (
	"gopkg.in/mgo.v2"
	"gitlab.com/iotTracker/brain/log"
	"gopkg.in/mgo.v2/bson"
)

type mongoRecordHandler struct{
	mongoSession *mgo.Session
	database, collection string
}

func NewMongoRecordHandler(mongoSession *mgo.Session, database, collection string) *mongoRecordHandler {

	setupRecords(mongoSession, database, collection)

	NewMongoRecordHandler := mongoRecordHandler{
		mongoSession,
		database,
		collection,
	}

	err := initialSystemRoleSetup(&NewMongoRecordHandler)
	if err != nil {
		log.Fatal("Unable to complete Initial System Role Setup!", err)
	}

	return &NewMongoRecordHandler
}

func setupRecords(mongoSession *mgo.Session, database, collection string){
	//Initialise record collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	roleCollection := mgoSesh.DB(database).C(collection)

	//Ensure index Uniqueness
	uniqueIndex := mgo.Index{
		Key: []string{"name"},
		Unique: true,
	}
	if err := roleCollection.EnsureIndex(uniqueIndex); err != nil {
		log.Fatal("Could not ensure uniqueness on name index in role collection: ", err)
	}

}

func initialSystemRoleSetup(handler *mongoRecordHandler) error {
	for _, roleToCreate := range initialRoles {
		//Try and retrieve the record
		retrieveRoleResponse := RetrieveResponse{}

		if err := handler.Retrieve(&RetrieveRequest{roleToCreate.Name}, &retrieveRoleResponse); err != nil {
			//Unable to retrieve role record
			//Try create role record
			createRoleResponse := CreateResponse{}

			if err := handler.Create(&CreateRequest{roleToCreate}, &createRoleResponse); err != nil {
				log.Fatal("Unable to Complete Initial Role Setup!", "Could Not Create role: " + roleToCreate.Name)
			}

			log.Info("Initial Role Setup: Created Role: " + roleToCreate.Name)
			continue
		}

		//Record Retrieved Successfully
		if roleToCreate.ComparePermissions(retrieveRoleResponse.SystemRole.Permissions) {
			log.Info("Initial Role Setup: Role " + retrieveRoleResponse.SystemRole.Name + " already exists and permissions correct.")
			continue
		}

		log.Info("Initial Role Setup: Role " + retrieveRoleResponse.SystemRole.Name + " already exists. Updating permissions.")
		updateRoleResponse := UpdateResponse{}
		if err := handler.Update(&UpdateRequest{roleToCreate}, &updateRoleResponse); err != nil {
			log.Fatal("Unable to Complete Initial Role Setup! Error updating role permissions!")
		}
	}

	return nil
}

func (u *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(u.database).C(u.collection)

	err := roleCollection.Insert(request.SystemRole)

	if err != nil {
		log.Error("Could not create Role! ", err)
		return err //TODO: Translate Unknown error
	}
	return nil
}

func (u *mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(u.database).C(u.collection)
	err := roleCollection.Find(bson.M{"name":request.Name}).One(&response.SystemRole)
	if err != nil {
		//log.Error("Unable to retrieve role!", err)
		return err
	}

	return nil
}


func (u *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(u.database).C(u.collection)

	err := roleCollection.Update(bson.M{"name":request.SystemRole.Name}, request.SystemRole)
	if err != nil {
		log.Error("Unable to update role!", err)
	}

	return nil
}
