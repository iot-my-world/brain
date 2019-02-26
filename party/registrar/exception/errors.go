package exception

import "strings"

type UnableToRetrieveParty struct {
	Reasons []string
}

func (e UnableToRetrieveParty) Error() string {
	return "unable to retrieve party: " + strings.Join(e.Reasons, "; ")
}

type AlreadyRegistered struct {}

func (e AlreadyRegistered) Error() string {
	return "user already registered"
}
