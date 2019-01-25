package businessDay

import (
	"gopkg.in/mgo.v2"
	"bitbucket.org/gotimekeeper/log"
	"time"
	"strings"
	"errors"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	businessDayConfig "bitbucket.org/gotimekeeper/business/businessDay/config"
	shiftConfig "bitbucket.org/gotimekeeper/business/shift/config"
	"bitbucket.org/gotimekeeper/business/shift"
	"bitbucket.org/gotimekeeper/business/businessRole"
	"bitbucket.org/gotimekeeper/business"
)

type mongoRecordHandler struct{
	mongoSession                   *mgo.Session
	database, collection           string
	businessDayConfigRecordHandler businessDayConfig.RecordHandler
	businessRoleRecordHandler      businessRole.RecordHandler
	employeeCollection             string
}

func NewMongoRecordHandler(
	mongoSession 					*mgo.Session,
	database, collection 			string,
	businessDayConfigRecordHandler 	businessDayConfig.RecordHandler,
	businessRoleRecordHandler 		businessRole.RecordHandler,
	employeeCollection             string,
	) *mongoRecordHandler {

	setupRecords(mongoSession, database, collection)

	NewMongoRecordHandler := mongoRecordHandler{
		mongoSession:                   mongoSession,
		database:                       database,
		businessDayConfigRecordHandler: businessDayConfigRecordHandler,
		businessRoleRecordHandler:      businessRoleRecordHandler,
		employeeCollection:          	employeeCollection,
		collection:                     collection,
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
		Key: []string{"id"},
		Unique: true,
	}
	if err := roleCollection.EnsureIndex(uniqueIndex); err != nil {
		log.Fatal("Could not ensure uniqueness on id in businessRole collection: ", err)
	}
}

func (m *mongoRecordHandler) validateReq(r interface{}) (error) {
	var reasonsInvalid []string

	switch v := r.(type){
	case *CreateRequest:
	case *GetCurrentRequest:
		deviationUpperBound := time.Now().Unix() + 1 * 60 * 60 // 1 Hour in Seconds
		deviationLowerBound := time.Now().Unix() - 1 * 60 * 60
		if v.ClientDateTime > deviationUpperBound ||
			v.ClientDateTime < deviationLowerBound {
			reasonsInvalid = append(
				reasonsInvalid,
				"Client Date and Time deviates from system by more than an hour",
				)
		}
	case *GetBeforeRequest:
	case *GetAfterRequest:
	case *GetSelectedRequest:
	case *UpdateShiftsRequest:
	case *EmployeeClockRequest:
	case *AssignEmployeesToShiftRequest:
		if v.BusinessDay.Id == "" {
			reasonsInvalid = append(
				reasonsInvalid,
				"businessDay Id Cannot Be Blank",
			)
		}
	default:
		log.Warn("NO CHECK CASE FOR THIS REQUEST!")
	}

	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, ","))
	}
	return nil
}

func (m *mongoRecordHandler) Create(request *CreateRequest, response *CreateResponse) error {
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	request.BusinessDay.Id = bson.NewObjectId().Hex()

	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()

	roleCollection := mgoSession.DB(m.database).C(m.collection)

	err = roleCollection.Insert(request.BusinessDay)
	if err != nil {
		log.Error("Could not create business day", err)
		return err //TODO: Translate Unknown error
	}

	response.BusinessDay = request.BusinessDay

	return nil
}

var months = map[string] string {
	"January":   "01",
	"February":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

func (m *mongoRecordHandler) GetCurrent(request *GetCurrentRequest, response *GetCurrentResponse) error {
	// Validate request data
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	// Get the start and end of today in Unix
	year :=  strconv.Itoa(time.Now().Year())
	month := months[time.Now().Month().String()]
	day := strconv.Itoa(time.Now().Day())
	if time.Now().Day() < 10 {
		day = "0" + day
	}

	startOfTodayStr := year + "-" + month + "-" + day +"T00:00:00+00:00"
	endOfTodayStr := year + "-" + month + "-" + day +"T23:59:59+00:00"
	startOfTodayTime, err := time.Parse(time.RFC3339, startOfTodayStr)
	if err != nil {
		log.Error("Unable to parse start of day")
		return errors.New("unable to find day range")
	}

	endOfTodayTime, err := time.Parse(time.RFC3339, endOfTodayStr)
	if err != nil {
		log.Error("Unable to parse end of day")
		return errors.New("unable to find day range")
	}

	startOfTodayUnix := startOfTodayTime.Unix()
	endOfTodayUnix := endOfTodayTime.Unix()

	// Build filter for dateTimeRange of today
	filter := make(map[string]interface{})
	filter["startDateTime"] = map[string]int64{"$gte": startOfTodayUnix, "$lte": endOfTodayUnix}

	// Open mongo session
	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(m.database).C(m.collection)

	var records []business.BusinessDay
	if err := collection.Find(filter).All(&records); err != nil {
		log.Warn("unable to retrieve BusinessDay records!")
		return err
	}

	if len(records) > 1 {
		// More than one record found...this should never happen.
		// There should never be more than one business day starting on the same day
		log.Warn("2 Business Day Records found. This should never happen.")
		response.BusinessDay = records[0]
		return nil
	} else if len(records) == 1 {
		response.BusinessDay = records[0]
		return nil
	}

	// len(records) == 0
	// We must create a new one for today
	// Retrieve businessDayConfig For Today
	latestBusinessDayRetrieveResponse := businessDayConfig.RetrieveResponse{}
	m.businessDayConfigRecordHandler.Retrieve(
		&businessDayConfig.RetrieveRequest{},
		&latestBusinessDayRetrieveResponse,
	)

	todaysConfig := map[string][]shiftConfig.Config {
		"Monday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Monday,
		"Tuesday":   latestBusinessDayRetrieveResponse.BusinessDayConfig.Tuesday,
		"Wednesday": latestBusinessDayRetrieveResponse.BusinessDayConfig.Wednesday,
		"Thursday":  latestBusinessDayRetrieveResponse.BusinessDayConfig.Thursday,
		"Friday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Friday,
		"Saturday":  latestBusinessDayRetrieveResponse.BusinessDayConfig.Saturday,
		"Sunday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Sunday,
	}[startOfTodayTime.Weekday().String()]

	// Make new business day
	newBusinessDay := business.BusinessDay{}

	// Confirm that 1 or more shifts are configured for business day
	if len(todaysConfig) == 0 {
		return errors.New("no shifts configured for business day " + startOfTodayTime.Weekday().String())
	}

	// Having retrieved the config,
	// We now add business roles to this business day we are creating
	totalNoShiftHoursInDay := 0
	for shiftIndex, shiftConfigEntry := range todaysConfig {
		// Create new shift
		newShift := shift.Shift{}

		// For Shift StartDateTime:
		// Get DateTime in go time structure
		newShift.StartDateTime = shiftConfigEntry.StartDateTime
		newShiftStartTime := time.Unix(newShift.StartDateTime, 0)
		// Increment year, month and day to reflect server dateTime
		newShiftStartTime = newShiftStartTime.AddDate(
			startOfTodayTime.Year() - newShiftStartTime.Year(),
			int(startOfTodayTime.Month()) - int(newShiftStartTime.Month()),
			startOfTodayTime.Day() - newShiftStartTime.Day(),
			)
		// Set newShift DateTime to unix of result after setting year, month and day
		newShift.StartDateTime = newShiftStartTime.Unix()

		// For Shift StartDateTime:
		// Get DateTime in go time structure
		newShift.EndDateTime = shiftConfigEntry.EndDateTime
		newShiftEndTime := time.Unix(newShift.EndDateTime, 0)
		// Increment year, month and day to reflect server dateTime
		newShiftEndTime = newShiftEndTime.AddDate(
			startOfTodayTime.Year() - newShiftEndTime.Year(),
			int(startOfTodayTime.Month()) - int(newShiftEndTime.Month()),
			startOfTodayTime.Day() - newShiftEndTime.Day(),
		)
		// Set newShift DateTime to unix of result after setting year, month and day
		newShift.EndDateTime = newShiftEndTime.Unix()

		// Get length of this shift in hours
		shiftLengthHours := newShiftEndTime.Hour() - newShiftStartTime.Hour()
		if shiftLengthHours < 0 {
			// If difference is negative add 24
			// (This shift will be ending tomorrow sometime)
			shiftLengthHours += 24
		}

		// Increment total no. of shift hours accumulated this day so fay
		totalNoShiftHoursInDay += shiftLengthHours

		// If it has become equal to 24, we need to add a day to our newShiftEndTime
		if totalNoShiftHoursInDay == 24 {
			newShiftEndTime = newShiftEndTime.AddDate(0, 0, 1)
			newShift.EndDateTime = newShiftEndTime.Unix()
		} else if totalNoShiftHoursInDay > 24 {
			return errors.New("total no of hours configured for business day " + startOfTodayTime.Weekday().String()  + " exceeds 24")
		}
		if shiftIndex == len(todaysConfig) - 1 {
			if totalNoShiftHoursInDay < 24 {
				log.Warn("total no of hours configured for business day " + startOfTodayTime.Weekday().String()  + " is less than 24")
			}
		}
		// Give shift new Id
		newShift.Id = bson.NewObjectId().Hex()
		// Add shift to this business day's shifts
		newBusinessDay.Shifts = append(newBusinessDay.Shifts, newShift)
	}

	// Set start and end dateTime of business day by start and end dateTime of first and last shifts respectively
	newBusinessDay.StartDateTime = newBusinessDay.Shifts[0].StartDateTime
	newBusinessDay.EndDateTime = newBusinessDay.Shifts[len(newBusinessDay.Shifts) - 1].EndDateTime

	createBusinessDayResponse := CreateResponse{}
	if err := m.Create(&CreateRequest{newBusinessDay}, &createBusinessDayResponse); err != nil {
		log.Error("Unable to create business day")
		return errors.New("unable to create new business day")
	}

	response.BusinessDay = createBusinessDayResponse.BusinessDay

	return nil
}

func (m *mongoRecordHandler) GetBefore(request *GetBeforeRequest, response *GetBeforeResponse) error {
	// Validate request data
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	// Build filter for dateTimeRange of the business day in request
	filter := make(map[string]interface{})
	// We want to find a business day record which has a start date time within 24 hours before
	// the startDateTime of the business day given in the request
	filter["startDateTime"] = map[string]int64{
		"$gte": request.BusinessDay.StartDateTime - 1 * 60 * 60 * 24, // 24 hours in seconds
		"$lt": request.BusinessDay.StartDateTime,
	}

	// Open mongo session
	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(m.database).C(m.collection)

	var records []business.BusinessDay
	if err := collection.Find(filter).All(&records); err != nil {
		log.Warn("unable to retrieve BusinessDay records!")
		return err
	}

	if len(records) > 1 {
		// More than one record found...this should never happen.
		// There should never be more than one business day starting on the same day
		log.Warn("2 Business Day Records found. This should never happen.")
		response.BusinessDay = records[0]
		return nil
	} else if len(records) == 1 {
		response.BusinessDay = records[0]
		return nil
	}

	// len(records) == 0
	// TODO: Return here? Throw previous business day does not exist error? Or Create?
	return errors.New("business day before does not exist")


	//// We must create a new one the day in the future being fetched
	//// Retrieve the latest business day config
	//latestBusinessDayRetrieveResponse := businessDayConfig.RetrieveResponse{}
	//m.businessDayConfigRecordHandler.Retrieve(
	//	&businessDayConfig.RetrieveRequest{},
	//	&latestBusinessDayRetrieveResponse,
	//)
	//
	//
	//requestBusinessDayStartDateTime := time.Unix(request.StartDateTime, 0)
	//responseBusinessDayOfWeek := ""
	//if requestBusinessDayStartDateTime.Weekday().String() == "Monday" {
	//	responseBusinessDayOfWeek = "Sunday"
	//} else {
	//	responseBusinessDayOfWeek = time.Unix(
	//		request.StartDateTime - 1 * 60 * 60 * 24, // 24 hours in seconds
	//		0).Weekday().String()
	//}
	//
	//responseBusinessDayConfig := map[string][]shiftConfig.Config {
	//	"Monday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Monday,
	//	"Tuesday":   latestBusinessDayRetrieveResponse.BusinessDayConfig.Tuesday,
	//	"Wednesday": latestBusinessDayRetrieveResponse.BusinessDayConfig.Wednesday,
	//	"Thursday":  latestBusinessDayRetrieveResponse.BusinessDayConfig.Thursday,
	//	"Friday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Friday,
	//	"Saturday":  latestBusinessDayRetrieveResponse.BusinessDayConfig.Saturday,
	//	"Sunday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Sunday,
	//}[responseBusinessDayOfWeek]

	//return nil
}

func (m *mongoRecordHandler) ConstructNewBusinessDayRecord(assignedShifts []shiftConfig.Config, startOfTodayTime time.Time) (business.BusinessDay, error) {
	// Make new business day
	newBusinessDay := business.BusinessDay{}

	if len(assignedShifts) == 0 {
		return business.BusinessDay{}, errors.New("no shifts configured for business day")
	}

	// Confirm that 1 or more shifts are configured for business day
	if len(assignedShifts) == 0 {
		return business.BusinessDay{}, errors.New("no shifts configured for business day")
	}

	// Having retrieved the config,
	// We now add business roles to this business day we are creating
	totalNoShiftHoursInDay := 0
	for shiftIndex, shiftConfigEntry := range assignedShifts {
		// Create new shift
		newShift := shift.Shift{}

		// For Shift StartDateTime:
		// Get DateTime in go time structure
		newShift.StartDateTime = shiftConfigEntry.StartDateTime
		newShiftStartTime := time.Unix(newShift.StartDateTime, 0)
		// Increment year, month and day to reflect server dateTime
		newShiftStartTime = newShiftStartTime.AddDate(
			startOfTodayTime.Year() - newShiftStartTime.Year(),
			int(startOfTodayTime.Month()) - int(newShiftStartTime.Month()),
			startOfTodayTime.Day() - newShiftStartTime.Day(),
		)
		// Set newShift DateTime to unix of result after setting year, month and day
		newShift.StartDateTime = newShiftStartTime.Unix()

		// For Shift StartDateTime:
		// Get DateTime in go time structure
		newShift.EndDateTime = shiftConfigEntry.EndDateTime
		newShiftEndTime := time.Unix(newShift.EndDateTime, 0)
		// Increment year, month and day to reflect server dateTime
		newShiftEndTime = newShiftEndTime.AddDate(
			startOfTodayTime.Year() - newShiftEndTime.Year(),
			int(startOfTodayTime.Month()) - int(newShiftEndTime.Month()),
			startOfTodayTime.Day() - newShiftEndTime.Day(),
		)
		// Set newShift DateTime to unix of result after setting year, month and day
		newShift.EndDateTime = newShiftEndTime.Unix()

		// Get length of this shift in hours
		shiftLengthHours := newShiftEndTime.Hour() - newShiftStartTime.Hour()
		if shiftLengthHours < 0 {
			// If difference is negative add 24
			// (This shift will be ending tomorrow sometime)
			shiftLengthHours += 24
		}

		// Increment total no. of shift hours accumulated this day so fay
		totalNoShiftHoursInDay += shiftLengthHours

		// If it has become equal to 24, we need to add a day to our newShiftEndTime
		if totalNoShiftHoursInDay == 24 {
			newShiftEndTime = newShiftEndTime.AddDate(0, 0, 1)
			newShift.EndDateTime = newShiftEndTime.Unix()
		} else if totalNoShiftHoursInDay > 24 {
			return business.BusinessDay{}, errors.New("total no of hours configured for business day " + startOfTodayTime.Weekday().String()  + " exceeds 24")
		}
		if shiftIndex == len(assignedShifts) - 1 {
			if totalNoShiftHoursInDay < 24 {
				log.Warn("total no of hours configured for business day " + startOfTodayTime.Weekday().String()  + " is less than 24")
			}
		}
		// Give shift new Id
		newShift.Id = bson.NewObjectId().Hex()
		// Add shift to this business day's shifts
		newBusinessDay.Shifts = append(newBusinessDay.Shifts, newShift)
	}

	// Set start and end dateTime of business day by start and end dateTime of first and last shifts respectively
	newBusinessDay.StartDateTime = newBusinessDay.Shifts[0].StartDateTime
	newBusinessDay.EndDateTime = newBusinessDay.Shifts[len(newBusinessDay.Shifts) - 1].EndDateTime

	return newBusinessDay, nil
}

func (m *mongoRecordHandler) GetAfter(request *GetAfterRequest, response *GetAfterResponse) error {
	// Validate request data
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	// Build filter for dateTimeRange of the business day in request
	filter := make(map[string]interface{})
	// We want to find a business day record which has a start date time within 24 hours of the
	// business day given in the request
	filter["startDateTime"] = map[string]int64{
		"$gte": request.BusinessDay.EndDateTime,
		"$lt": request.BusinessDay.EndDateTime + 1 * 60 * 60 * 24, // 24 hours in seconds
	}

	// Open mongo session
	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(m.database).C(m.collection)

	var records []business.BusinessDay
	if err := collection.Find(filter).All(&records); err != nil {
		log.Warn("unable to retrieve BusinessDay records!")
		return err
	}

	if len(records) > 1 {
		// More than one record found...this should never happen.
		// There should never be more than one business day starting on the same day
		log.Warn("2 Business Day Records found. This should never happen.")
		response.BusinessDay = records[0]
		return nil
	} else if len(records) == 1 {
		response.BusinessDay = records[0]
		return nil
	}

	// len(records) == 0
	// We must create a new one the day in the future being fetched
	// Retrieve the latest business day config
	latestBusinessDayRetrieveResponse := businessDayConfig.RetrieveResponse{}
	m.businessDayConfigRecordHandler.Retrieve(
		&businessDayConfig.RetrieveRequest{},
		&latestBusinessDayRetrieveResponse,
	)


	requestBusinessDayStartDateTime := time.Unix(request.StartDateTime, 0)
	responseBusinessDayOfWeek := ""
	if requestBusinessDayStartDateTime.Weekday().String() == "Sunday" {
		responseBusinessDayOfWeek = "Monday"
	} else {
		responseBusinessDayOfWeek = time.Unix(
			request.StartDateTime + 1 * 60 * 60 * 24, // 24 hours in seconds
			0).Weekday().String()
	}

	responseBusinessStartDateTime := time.Unix(
		request.StartDateTime + 1 * 60 * 60 * 24, // 24 hours in seconds
		0)

	responseBusinessDayConfig := map[string][]shiftConfig.Config {
		"Monday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Monday,
		"Tuesday":   latestBusinessDayRetrieveResponse.BusinessDayConfig.Tuesday,
		"Wednesday": latestBusinessDayRetrieveResponse.BusinessDayConfig.Wednesday,
		"Thursday":  latestBusinessDayRetrieveResponse.BusinessDayConfig.Thursday,
		"Friday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Friday,
		"Saturday":  latestBusinessDayRetrieveResponse.BusinessDayConfig.Saturday,
		"Sunday":    latestBusinessDayRetrieveResponse.BusinessDayConfig.Sunday,
	}[responseBusinessDayOfWeek]

	newBusinessDay := business.BusinessDay{}
	if newBusinessDay, err = m.ConstructNewBusinessDayRecord(responseBusinessDayConfig, responseBusinessStartDateTime); err != nil {
		return errors.New("unable to create new business day")
	}

	createBusinessDayResponse := CreateResponse{}
	if err := m.Create(&CreateRequest{newBusinessDay}, &createBusinessDayResponse); err != nil {
		log.Error("Unable to create business day")
		return errors.New("unable to create new business day")
	}

	response.BusinessDay = createBusinessDayResponse.BusinessDay

	return nil
}

func (m *mongoRecordHandler) GetSelected(request *GetSelectedRequest, response *GetSelectedResponse) error {
	// Validate request data
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	// Open connection to database
	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(m.database).C(m.collection)

	existingRecord := business.BusinessDay{}

	// Try and retrieve existing record
	if err := collection.Find(bson.M{"id": request.BusinessDay.Id}).One(&existingRecord); err != nil {
		return err
	}

	// Set the retrieved record
	response.BusinessDay = existingRecord

	return nil
}

func (m *mongoRecordHandler) UpdateShifts(request *UpdateShiftsRequest, response *UpdateShiftsResponse) error {
	// Validate request data
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	// Open connection to database
	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(m.database).C(m.collection)

	existingRecord := business.BusinessDay{}

	// Try and retrieve existing record
	if err := collection.Find(bson.M{"id": request.BusinessDay.Id}).One(&existingRecord); err != nil {
		return err
	}

	// Update Shifts
	existingRecord.Shifts = request.BusinessDay.Shifts

	// Try insert update
	if err := collection.Update(bson.M{"id": request.BusinessDay.Id}, existingRecord); err != nil {
		return err
	}

	response.BusinessDay = existingRecord

	return nil
}

func (m *mongoRecordHandler) AssignEmployeesToShift(request *AssignEmployeesToShiftRequest, response *AssignEmployeesToShiftResponse) error {
	// Validate request data
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	// Open connection to database
	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()
	businessDayCollection := mgoSession.DB(m.database).C(m.collection)
	employeeCollection := mgoSession.DB(m.database).C(m.employeeCollection)

	businessDayRecord := business.BusinessDay{}

	// Try and retrieve existing business day record
	if err := businessDayCollection.Find(bson.M{"id": request.BusinessDay.Id}).One(&businessDayRecord); err != nil {
		return err
	}

	// Confirm that ids are of valid employees
	employeesToAdd := make([]business.Employee, 0)
	for _, employeeId := range request.EmployeeIds {
		employeeToAdd := business.Employee{}
		if err := employeeCollection.Find(bson.M{"id": employeeId}).One(&employeeToAdd); err == nil {
			employeesToAdd = append(employeesToAdd, employeeToAdd)
		}
	}

	// For Each Shift in the given business record...
	reInsert := false
	for i, shiftRecord := range businessDayRecord.Shifts {
		// Find the shift given id
		if request.ShiftId == shiftRecord.Id {
			for _, emp := range employeesToAdd {

				// Check if this employee already has a register entry
				createRegEnt := true
				for _, regEnt := range shiftRecord.Register {
					if regEnt.EmployeeId == emp.Id {
						createRegEnt = false
					}
				}

				if createRegEnt {
					reInsert = true
					businessDayRecord.Shifts[i].Register = append(
						businessDayRecord.Shifts[i].Register,
						shift.RegEntry{
							EmployeeId: emp.Id,
						},
						)
				}
			}
		}
	}

	if reInsert {
		// Try insert update
		if err := businessDayCollection.Update(bson.M{"id": request.BusinessDay.Id}, businessDayRecord); err != nil {
			return err
		}
	}

	response.BusinessDay = businessDayRecord

	return nil
}

func (m *mongoRecordHandler) EmployeeClock(request *EmployeeClockRequest, response *EmployeeClockResponse) error {
	// Validate request data
	err := m.validateReq(request)
	if err != nil {
		return err
	}

	// Get Current Business Day
	getCurrentBusienssDayResponse := GetCurrentResponse{}
	if err := m.GetCurrent(&GetCurrentRequest{ClientDateTime: time.Now().Unix()}, &getCurrentBusienssDayResponse); err != nil {
		log.Error("Error getting current business day: " + err.Error())
		return errors.New("Error getting current business day: " + err.Error())
	}
	currentBusinessDay := getCurrentBusienssDayResponse.BusinessDay


	// Check if employee is already here
	// (i.e in clocked in list of employee IDs)
	employeeIn := false
	for _, employeeID := range currentBusinessDay.ClockedIn {
		if employeeID == request.Employee.Id {
			employeeIn = true
			break
		}
	}
	if employeeIn {
		// Employee is in, Clock Them Out
		// Create Clock-Out Event in current business day clock history
		clockOut := business.ClockEvent{
			Direction: business.CLOCK_OUT,
			EmployeeId: request.Employee.Id,
		}
		response.ClockEvent = clockOut
		currentBusinessDay.ClockHistory = append(currentBusinessDay.ClockHistory, clockOut)

		// Add Employee ID to clocked-out list on current business day
		currentBusinessDay.ClockedOut = append(
			currentBusinessDay.ClockedOut,
			request.Employee.Id,
		)

		// Remove Employee ID From Clocked-in
		for empNo, employeeID := range currentBusinessDay.ClockedIn {
			if employeeID == request.Employee.Id {
				if len(currentBusinessDay.ClockedIn) == 1 {
					currentBusinessDay.ClockedIn = make([]string, 0)
				} else {
					// Remove employee id from clocked out list
					currentBusinessDay.ClockedIn = append(currentBusinessDay.ClockedIn[:empNo], currentBusinessDay.ClockedIn[empNo+1:]...)
				}
			}
		}
	} else {
		// Employee is out, clock them in
		// Create Clock-In Event in business day history
		clockIn := business.ClockEvent{
			Direction: business.CLOCK_IN,
			EmployeeId: request.Employee.Id,
		}
		response.ClockEvent = clockIn
		currentBusinessDay.ClockHistory = append(currentBusinessDay.ClockHistory, clockIn)


		// Remove employee's ID from clockedOut list if it is already there
		// (i.e. they were at work but left and are coming back)
		for empNo, employeeID := range currentBusinessDay.ClockedOut {
			if employeeID == request.Employee.Id {
				if len(currentBusinessDay.ClockedOut) == 1 {
					currentBusinessDay.ClockedOut = make([]string, 0)
				} else {
					// Remove employee id from clocked out list
					currentBusinessDay.ClockedOut = append(currentBusinessDay.ClockedOut[:empNo], currentBusinessDay.ClockedOut[empNo+1:]...)
				}

			}
		}

		// Add employee's id to clockedIn
		currentBusinessDay.ClockedIn = append(
			currentBusinessDay.ClockedIn,
			request.Employee.Id,
		)

	}

	//fmt.Println("Clocked In:", currentBusinessDay.ClockedIn)
	//fmt.Println("Clocked Out:", currentBusinessDay.ClockedOut)
	//fmt.Println("Clock History:")
	//for _, clockEvent := range currentBusinessDay.ClockHistory {
	//	fmt.Println(clockEvent)
	//}

	// Open connection to database
	mgoSession := m.mongoSession.Copy()
	defer mgoSession.Close()
	collection := mgoSession.DB(m.database).C(m.collection)

	// Save Updated Business Day
	// Try insert update
	if err := collection.Update(bson.M{"id": currentBusinessDay.Id}, currentBusinessDay); err != nil {
		log.Error("Error while updating current business day: "  + err.Error())
		return errors.New("Error while updating current business day: "  + err.Error())
	}

	// Put updated business day into the response
	response.BusinessDay = currentBusinessDay
	return nil
}