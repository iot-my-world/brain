package config

type Config struct {
	StartDateTime int64                      `json:"startDateTime" bson:"startDateTime"`
	EndDateTime   int64                      `json:"endDateTime" bson:"endDateTime"`
}
