package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error() string {
	return "initial tk102 setup error: " + strings.Join(e.Reasons, "; ")
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "tk102 not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "tk102 creation error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "tk102 update error: " + strings.Join(e.Reasons, "; ")
}
