package exception

import "strings"

type ClientCreation struct {
	Reasons []string
}

func (e ClientCreation) Error() string {
	return "client creation error: " + strings.Join(e.Reasons, "; ")
}

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

type Delete struct {
	Reasons []string
}

func (e Delete) Error() string {
	return "delete error: " + strings.Join(e.Reasons, "; ")
}
