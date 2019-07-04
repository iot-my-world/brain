package exception

import "strings"

type NotFound struct{}

func (e NotFound) Error() string {
	return "user not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "user creation error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

type Delete struct {
	Reasons []string
}

func (e Delete) Error() string {
	return "user delete error: " + strings.Join(e.Reasons, "; ")
}

func (e Update) Error() string {
	return "user update error: " + strings.Join(e.Reasons, "; ")
}

type RecordHandlerNil struct{}

func (e RecordHandlerNil) Error() string {
	return "given brain client recordHandler is nil"
}

type Collect struct {
	Reasons []string
}

func (e Collect) Error() string {
	return "user collect error: " + strings.Join(e.Reasons, "; ")
}
