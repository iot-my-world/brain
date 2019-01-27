package config

import shift "gitlab.com/iotTracker/brain/business/shift/config"

type Config struct {
	Id string `json:"id" bson:"id"`
	Monday    []shift.Config `json:"monday" bson:"monday"`
	Tuesday   []shift.Config `json:"tuesday" bson:"tuesday"`
	Wednesday []shift.Config `json:"wednesday" bson:"wednesday"`
	Thursday  []shift.Config `json:"thursday" bson:"thursday"`
	Friday    []shift.Config `json:"friday" bson:"friday"`
	Saturday  []shift.Config `json:"saturday" bson:"saturday"`
	Sunday    []shift.Config `json:"sunday" bson:"sunday"`
}