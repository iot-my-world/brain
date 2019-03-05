package invite

import "gitlab.com/iotTracker/brain/party/user"

type Invite struct {
	Id    string    `json:"id" bson:"id"`
	User  user.User `json:"user" bson:"user"`
	token string    `json:"token" bson:"token"`
}
