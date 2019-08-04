package exception

import "strings"

type RecordHandlerNil struct{}

func (e RecordHandlerNil) Error() string {
	return "given brain message recordHandler is nil"
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "message not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "message creation error: " + strings.Join(e.Reasons, "; ")
}

type Retrieve struct {
	Reasons []string
}

func (e Retrieve) Error() string {
	return "message retrieval error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "message update error: " + strings.Join(e.Reasons, "; ")
}

type Delete struct {
	Reasons []string
}

func (e Delete) Error() string {
	return "message delete error: " + strings.Join(e.Reasons, "; ")
}

type Collect struct {
	Reasons []string
}

func (e Collect) Error() string {
	return "message collect error: " + strings.Join(e.Reasons, "; ")
}
