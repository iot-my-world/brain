package exception

import "strings"

type NotFound struct{}

func (e NotFound) Error() string {
	return "system not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "system creation error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "system update error: " + strings.Join(e.Reasons, "; ")
}
