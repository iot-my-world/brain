package business


type Employee struct {
	Id               string `json:"id" bson:"id"`
	Name             string `json:"name" bson:"name"`
	Surname          string `json:"surname" bson:"surname"`
	IDNo             uint64 `json:"idNo" bson:"idNo"`
	TagID            string `json:"tagID" bson:"tagID"`
	Picture          string `json:"picture" bson:"picture"`
	ContactNo        string `json:"contactNo" bson:"contactNo"`
	MondayShiftNo    int    `json:"mondayShiftNo" bson:"mondayShiftNo"`
	TuesdayShiftNo   int    `json:"tuesdayShiftNo" bson:"tuesdayShiftNo"`
	WednesdayShiftNo int    `json:"wednesdayShiftNo" bson:"wednesdayShiftNo"`
	ThursdayShiftNo  int    `json:"thursdayShiftNo" bson:"thursdayShiftNo"`
	FridayShiftNo    int    `json:"fridayShiftNo" bson:"fridayShiftNo"`
	SaturdayShiftNo  int    `json:"saturdayShiftNo" bson:"saturdayShiftNo"`
	SundayShiftNo    int    `json:"sundayShiftNo" bson:"sundayShiftNo"`
}

