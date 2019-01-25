package config

type Config struct {
	BusinessRoleName string `json:"businessRoleName" bson:"businessRole"`
	BusinessRoleId   string `json:"businessRoleId" bson:"businessRoleId"`
	No               int    `json:"no" bson:"no"`
}