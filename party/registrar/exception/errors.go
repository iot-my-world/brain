package exception

import "strings"

type UnableToRetrieveParty struct {
	Reasons []string
}

func (e UnableToRetrieveParty) Error() string {
	return "unable to retrieve party: " + strings.Join(e.Reasons, "; ")
}

type AlreadyRegistered struct{}

func (e AlreadyRegistered) Error() string {
	return "user already registered"
}

type TokenGeneration struct {
	Reasons []string
}

func (e TokenGeneration) Error() string {
	return "token generation error: " + strings.Join(e.Reasons, "; ")
}
