package employee

import (
	"gopkg.in/mgo.v2"
	"bitbucket.org/gotimekeeper/log"
	"strings"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"bitbucket.org/gotimekeeper/business/businessDay"
	"time"
	"bitbucket.org/gotimekeeper/business"
)


type mongoRecordHandler struct{
	mongoSession *mgo.Session
	database, collection string
	businessDayRecordHandler businessDay.RecordHandler
}

func NewMongoRecordHandler(
	mongoSession *mgo.Session,
	database, collection string,
	businessDayRecordHandler businessDay.RecordHandler,
	) *mongoRecordHandler {

	setupRecords(mongoSession, database, collection)

	newMongoRecordHandler := mongoRecordHandler{
		mongoSession: mongoSession,
		database: database,
		collection: collection,
		businessDayRecordHandler: businessDayRecordHandler,
	}

	return &newMongoRecordHandler
}

func setupRecords(mongoSession *mgo.Session, database, collection string){
	//Initialise User collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	mgoCollection := mgoSesh.DB(database).C(collection)

	//Ensure system id Uniqueness
	idUnique := mgo.Index{
		Key: []string{"id"},
		Unique: true,
	}
	if err := mgoCollection.EnsureIndex(idUnique); err != nil {
		log.Fatal("Could not ensure uniqueness: ", err)
	}

	//Ensure Employee ID No Uniqueness
	idNoUnique := mgo.Index{
		Key: []string{"idNo"},
		Unique: true,
	}
	if err := mgoCollection.EnsureIndex(idNoUnique); err != nil {
		log.Fatal("Could not ensure uniqueness: ", err)
	}
}

func validateServiceReqData(r interface{}) error {
	var reasonsInvalid []string
	switch v := r.(type){
	case *CreateRequest:
		if v.Employee.Name == "" {
			reasonsInvalid = append(reasonsInvalid, "name cannot be blank")
		}

		//idNoLen := len(strconv.FormatUint(v.Employee.IDNo, 10))
		//if idNoLen != 13  {
		//	reasonsInvalid = append(reasonsInvalid, "ID Number must be 13 digits")
		//}
		// TODO: add more validation
	case *RetrieveAllRequest:
		//No validation required
	case *UpdateRequest:
		if v.Employee.Id == "" {
			reasonsInvalid = append(reasonsInvalid, "system id cannot be blank")
		}

	case *RetrieveByShiftAssignmentRequest:
	case *RetrieveByTagIDRequest:
		if v.TagID == "" {
			reasonsInvalid = append(reasonsInvalid, "tag ID cannot be blank when retrieving by tag id")
		}

	default:
		log.Warn("NO CHECK CASE FOR THIS REQUEST!")
	}

	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, ","))
	}
	return nil
}

func (u *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(u.database).C(u.collection)

	request.Employee.Id = bson.NewObjectId().Hex()

	if err := collection.Insert(&request.Employee); err != nil {
		log.Error("Could not create user! ", err)
		return err //TODO: Translate Unknown error
	}

	response.Employee = request.Employee

	return nil
}

func (u *mongoRecordHandler) RetrieveAll(request *RetrieveAllRequest, response *RetrieveAllResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	collection := mgoSession.DB(u.database).C(u.collection)

	var records []business.Employee

	if err := collection.Find(bson.M{}).All(&records); err != nil {
		return err
	}

	response.Records = records
	return nil
}

func updateAllowedFields (recordToUpdate *business.Employee, requestRecord business.Employee) {
	recordToUpdate.Name = requestRecord.Name
	recordToUpdate.Surname = requestRecord.Surname
	recordToUpdate.IDNo = requestRecord.IDNo
	recordToUpdate.ContactNo = requestRecord.ContactNo

	recordToUpdate.MondayShiftNo = requestRecord.MondayShiftNo
	recordToUpdate.TuesdayShiftNo = requestRecord.TuesdayShiftNo
	recordToUpdate.WednesdayShiftNo = requestRecord.WednesdayShiftNo
	recordToUpdate.ThursdayShiftNo = requestRecord.ThursdayShiftNo
	recordToUpdate.FridayShiftNo = requestRecord.FridayShiftNo
	recordToUpdate.SaturdayShiftNo = requestRecord.SaturdayShiftNo
	recordToUpdate.SundayShiftNo = requestRecord.SundayShiftNo
}

func (u *mongoRecordHandler) Update(request *UpdateRequest, response *UpdateResponse) error {
	// Validate Request
	if err := validateServiceReqData(request); err != nil {
		return err
	}

	// Open Mongo session and get collection
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(u.database).C(u.collection)

	//Try and retrieve existing record
	existingRecord := business.Employee{}
	if err := collection.Find(bson.M{"id": request.Employee.Id}).One(&existingRecord); err != nil {
		return err
	}

	// Update allowed fields
	updateAllowedFields(&existingRecord, request.Employee)

	// Try insert update
	if err := collection.Update(bson.M{"id": request.Employee.Id}, existingRecord); err != nil {
		return err
	}

	response.Employee = existingRecord

	return nil
}

func (u *mongoRecordHandler) RetrieveByShiftAssignment(request *RetrieveByShiftAssignmentRequest, response *RetrieveByShiftAssignmentResponse) error {
	// Validate Request
	if err := validateServiceReqData(request); err != nil {
		return err
	}

	// Try and retrieve business day
	getSelectedBusinessDayResponse := businessDay.GetSelectedResponse{}
	if err := u.businessDayRecordHandler.GetSelected(
		&businessDay.GetSelectedRequest{
			BusinessDay: request.BusinessDay,
		},
		&getSelectedBusinessDayResponse,
		); err != nil {
		return err
	}

	// Open connection to database
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(u.database).C(u.collection)

	dayOfWeek := strings.ToLower(time.Unix(getSelectedBusinessDayResponse.BusinessDay.StartDateTime, 0).Weekday().String())

	// Get all employees with an assigned role for this business day
	shiftGroups := make([][]business.Employee, 0)
	for shiftNo := range getSelectedBusinessDayResponse.BusinessDay.Shifts {
		shiftGroup := make([]business.Employee, 0)
		if err := collection.Find(bson.M{dayOfWeek + "ShiftNo" : shiftNo + 1}).All(&shiftGroup); err == nil {
			shiftGroups = append(shiftGroups, shiftGroup)
		}
	}

	// Get all unassigned employees
	unassigned := make([]business.Employee, 0)
	collection.Find(bson.M{dayOfWeek + "ShiftNo" : 0}).All(&unassigned)
	for _, unassignedEmployee := range unassigned {
		response.Unassigned = append(response.Unassigned, unassignedEmployee)
	}

	filter := make(map[string]interface{})
	filter["dayOfWeek"] = map[string]int{
		"$gt": len(getSelectedBusinessDayResponse.BusinessDay.Shifts),
	}
	collection.Find(filter).All(&unassigned)
	for _, unassignedEmployee := range unassigned {
		response.Unassigned = append(response.Unassigned, unassignedEmployee)
	}

	response.ShiftGroups = shiftGroups

	return nil
}

func (u *mongoRecordHandler) RetrieveByTagID(request *RetrieveByTagIDRequest, response *RetrieveByTagIDResponse) error {
	// Validate Request
	if err := validateServiceReqData(request); err != nil {
		return err
	}

	// Open Mongo session and get collection
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(u.database).C(u.collection)

	// Try and retrieve the employee
	existingRecord := business.Employee{}
	if err := collection.Find(bson.M{"tagID": request.TagID}).One(&existingRecord); err != nil {
		log.Error("Unable to retrieve employee with tag ID: " + request.TagID + " :" + err.Error())
		return errors.New("Unable to retrieve employee with tag ID: " + request.TagID + " :" + err.Error())
	}

	// If retrieve was successful, respond
	response.Employee = existingRecord

	return nil
}