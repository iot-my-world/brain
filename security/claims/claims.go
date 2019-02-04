package claims

import (
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type LoginClaims struct {
	UserId         id.Identifier `json:"userId"`
	IssuedAtTime   int64         `json:"issueTime"`
	ExpirationTime int64         `json:"expiry"`
}
