package exception

import "strings"

type RecordHandlerNil struct{}

func (e RecordHandlerNil) Error() string {
	return "given brain zx303 Task Reading recordHandler is nil"
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "zx303 Task Reading not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "zx303 Task Reading creation error: " + strings.Join(e.Reasons, "; ")
}

type Retrieve struct {
	Reasons []string
}

func (e Retrieve) Error() string {
	return "zx303 Task Reading retrieval error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "zx303 Task Reading update error: " + strings.Join(e.Reasons, "; ")
}

type Delete struct {
	Reasons []string
}

func (e Delete) Error() string {
	return "zx303 Task Reading delete error: " + strings.Join(e.Reasons, "; ")
}

type Collect struct {
	Reasons []string
}

func (e Collect) Error() string {
	return "zx303 Task Reading collect error: " + strings.Join(e.Reasons, "; ")
}
