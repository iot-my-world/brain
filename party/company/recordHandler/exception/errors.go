package exception

import "strings"

type RecordHandlerNil struct{}

func (e RecordHandlerNil) Error() string {
	return "given brain company recordHandler is nil"
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "company not found"
}

type Retrieval struct {
	Reasons []string
}

func (e Retrieval) Error() string {
	return "company retrieval error: " + strings.Join(e.Reasons)
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "company creation error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "company update error: " + strings.Join(e.Reasons, "; ")
}
