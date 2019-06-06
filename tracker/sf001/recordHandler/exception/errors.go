package exception

import "strings"

type RecordHandlerNil struct{}

func (e RecordHandlerNil) Error() string {
	return "given brain sf001 recordHandler is nil"
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "sf001 not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "sf001 creation error: " + strings.Join(e.Reasons, "; ")
}

type Retrieve struct {
	Reasons []string
}

func (e Retrieve) Error() string {
	return "sf001 retrieval error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "sf001 update error: " + strings.Join(e.Reasons, "; ")
}

type Delete struct {
	Reasons []string
}

func (e Delete) Error() string {
	return "sf001 delete error: " + strings.Join(e.Reasons, "; ")
}

type Collect struct {
	Reasons []string
}

func (e Collect) Error() string {
	return "sf001 collect error: " + strings.Join(e.Reasons, "; ")
}
