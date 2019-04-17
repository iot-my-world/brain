package exception

import "strings"

type RecordHandlerNil struct{}

func (e RecordHandlerNil) Error() string {
	return "given brain api user recordHandler is nil"
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "api user not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "api user creation error: " + strings.Join(e.Reasons, "; ")
}

type Retrieve struct {
	Reasons []string
}

func (e Retrieve) Error() string {
	return "api user retrieval error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "api user update error: " + strings.Join(e.Reasons, "; ")
}

type Delete struct {
	Reasons []string
}

func (e Delete) Error() string {
	return "api user delete error: " + strings.Join(e.Reasons, "; ")
}

type Collect struct {
	Reasons []string
}

func (e Collect) Error() string {
	return "api user collect error: " + strings.Join(e.Reasons, "; ")
}
