package exception

import "strings"

type NotFound struct{}

func (e NotFound) Error() string {
	return "not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "creation error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "update error: " + strings.Join(e.Reasons, "; ")
}
