package exception

import "strings"

type ClientRetrieval struct {
	Reasons []string
}

func (e ClientRetrieval) Error() string {
	return "client retrieval error: " + strings.Join(e.Reasons, "; ")
}

type AllowedFieldsUpdate struct {
	Reasons []string
}

func (e AllowedFieldsUpdate) Error() string {
	return "allowed fields update error: " + strings.Join(e.Reasons, "; ")
}
