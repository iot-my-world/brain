package tagEvent

import (
	"gopkg.in/mgo.v2"
	"bitbucket.org/gotimekeeper/log"
	"errors"
	"strings"
	"bitbucket.org/gotimekeeper/business/employee"
	"bitbucket.org/gotimekeeper/exoWSC"
	"bitbucket.org/gotimekeeper/exoWSC/message"
	"encoding/json"
	"bitbucket.org/gotimekeeper/rfId"
	"time"
	"fmt"
	"bitbucket.org/gotimekeeper/business/businessDay"
)

type mongoRecordHandler struct{
	mongoSession             *mgo.Session
	database, collection     string
	tagEventHub              *exoWSC.Hub
	employeeRecordHandler    employee.RecordHandler
	businessDayRecordHandler businessDay.RecordHandler
	MsgFromHub                chan exoWSC.Message
}

func NewMongoRecordHandler(
	mongoSession *mgo.Session,
	database, collection string,
	tagEventHub *exoWSC.Hub,
	employeeRecordHandler employee.RecordHandler,
	businessDayRecordHandler businessDay.RecordHandler,
	) *mongoRecordHandler {

	setupRecords(mongoSession, database, collection)

	NewMongoRecordHandler := mongoRecordHandler{
		mongoSession: mongoSession,
		database: database,
		collection: collection,
		tagEventHub: tagEventHub,
		employeeRecordHandler: employeeRecordHandler,
		businessDayRecordHandler: businessDayRecordHandler,
		MsgFromHub: make(chan exoWSC.Message, 1000),
	}

	return &NewMongoRecordHandler
}

func setupRecords(mongoSession *mgo.Session, database, collection string){
	//Initialise record collection in database
	mgoSesh := mongoSession.Copy()
	defer mgoSesh.Close()
	//tagEventCollection := mgoSesh.DB(database).C(collection)

	////Ensure index Uniqueness
	//uniqueIndex := mgo.Index{
	//	Key: []string{"tagId"},
	//	Unique: true,
	//}
	//if err := tagEventCollection.EnsureIndex(uniqueIndex); err != nil {
	//	log.Fatal("Could not ensure uniqueness on name index in tagEvent collection: ", err)
	//}

}

func validateServiceReqData(r interface{}) error {
	var reasonsInvalid []string

	switch v := r.(type){
	case *EmployeeClockRequest:

	case *RFIDTagEventRequest:
		if v.TagEvent.TagId == "" {
			reasonsInvalid = append(reasonsInvalid, "tagEvent tagId cannot be blank for a tagEvent create")
		}
	case *RetrieveRequest:
		if v.TagId == "" {
			reasonsInvalid = append(reasonsInvalid, "tagId cannot be blank for a tagEvent retrieve")
		}
	default:
		log.Warn("NO CHECK CASE FOR THIS REQUEST!")
	}
	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, ","))
	}
	return nil
}

func (u *mongoRecordHandler) Send(message exoWSC.Message) error {
	sendTimeOutTicker := time.NewTicker(2 *time.Second)
	defer func () {
		sendTimeOutTicker.Stop()
	}()

	select {
	case u.MsgFromHub <- message:
		break
	case <- sendTimeOutTicker.C:
		log.Error("Time out on waiting to get message into mongo tag event record handler's MsgFromHub Channel. Msg: ", message)
	}

	return nil
}

func (u *mongoRecordHandler) RFIDTagEvent(request *RFIDTagEventRequest, response *RFIDTagEventResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}
	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	newWSMsgByteData, err := json.Marshal(rfId.GetServiceContextRequest{
		Event: "tag",
	})
	if err != nil {
		log.Error("error marshalling rfid service context request: " + err.Error())
		return errors.New("error marshalling rfid service context request: " + err.Error())
	}
	messageData := string(newWSMsgByteData[:])

	newWSMsg := exoWSC.Message{
		Type: message.GetServiceContextRequest,
		SerialData: messageData,
	}


	timeOutTicker := time.NewTicker(1 * time.Second)
	defer func(){
		timeOutTicker.Stop()
	}()
	select {
	case u.tagEventHub.Broadcast <- newWSMsg:
		break
	case <- timeOutTicker.C:
		log.Error("time out waiting to put getServiceContextRequest message into hub's broadcast channel")
	}

	ticker := time.NewTicker(3 * time.Second)
	contextGot:
	for {
		select {
			case msg := <- u.MsgFromHub:
				switch msg.Type {
				case message.GetServiceContextResponse:
					contextResponse := rfId.GetServiceContextResponse{}
					if err := json.Unmarshal([]byte(msg.SerialData), &contextResponse); err != nil {
						log.Error("Error unmarshalling getServiceContextResponse for tag event: " + err.Error())
					}
					switch contextResponse.Context {
					case "assignToEmployee":
						fmt.Println("assign to employee")
					default:
						employeeClockResponse := EmployeeClockResponse{}
						if err := u.EmployeeClock(&EmployeeClockRequest{TagEvent:request.TagEvent}, &employeeClockResponse); err != nil {
							log.Error("Error during Employee Clock: " + err.Error())
							return errors.New("Error during Employee Clock: " + err.Error())
						}
					}

					break contextGot
				}

			case <- ticker.C:
				log.Error("time out while waiting for tag event context")
				break contextGot
		}
	}

	tagEventCollection := mgoSession.DB(u.database).C(u.collection)

	err = tagEventCollection.Insert(request.TagEvent)
	if err != nil {
		log.Error("Could not insert TagEvent ", err)
		return err
	}

	return nil
}

func (u *mongoRecordHandler) EmployeeClock(request *EmployeeClockRequest, response *EmployeeClockResponse) error {
	if err := validateServiceReqData(request); err != nil {
		return err
	}

	// Try and retrieve the employee that this tag belongs to
	employeeRetrieveByTagIDResponse := employee.RetrieveByTagIDResponse{}
	if err := u.employeeRecordHandler.RetrieveByTagID(
		&employee.RetrieveByTagIDRequest{
			TagID: request.TagEvent.TagId,
		},
		&employeeRetrieveByTagIDResponse); err != nil {
		log.Error("Error while retrieving employee by tag ID: " + err.Error())
		return errors.New("Error while retrieving employee by tag ID: " + err.Error())
	}

	// Call Employee Clock for updating of current business day
	businessDayEmployeeClockResponse := businessDay.EmployeeClockResponse{}
	if err := u.businessDayRecordHandler.EmployeeClock(
		&businessDay.EmployeeClockRequest{
			TagEvent: request.TagEvent,
			Employee: employeeRetrieveByTagIDResponse.Employee,
		},
		&businessDay.EmployeeClockResponse{},
		); err != nil {
		log.Error("Error while calling businessDay Employee Clock: " + err.Error())
		return errors.New("Error while calling businessDay Employee Clock: " + err.Error())
	}

	// Create Actual tag event

	mgoSession := u.mongoSession.Copy()
	defer mgoSession.Close()

	tagEventCollection := mgoSession.DB(u.database).C(u.collection)

	if err := tagEventCollection.Insert(request.TagEvent); err != nil {
		log.Error("Could not create TagEvent ", err)
		return errors.New("Could not create TagEvent " + err.Error())
	}

	// Create TagEventNotification and send it to be broadcast on the hub
	response.BusinessDay = businessDayEmployeeClockResponse.BusinessDay
	response.ClockEvent = businessDayEmployeeClockResponse.ClockEvent
	response.Employee = employeeRetrieveByTagIDResponse.Employee
	var serialResponse []byte
	var err error
	if serialResponse, err = json.Marshal(response); err != nil {
		log.Error("Error marshalling clock event response for WS: " + err.Error())
		return errors.New("Error marshalling clock event response for WS: " + err.Error())
	}

	tagEventNotificationWSMsg := exoWSC.Message{
		Type: message.ClockEvent,
		SerialData: string(serialResponse[:]),
	}
	u.tagEventHub.Broadcast <- tagEventNotificationWSMsg

	return nil
}

func (u *mongoRecordHandler) Retrieve(request *RetrieveRequest, response *RetrieveResponse) error {
	return nil
}