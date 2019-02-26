package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error() string {
	return "initial company setup error: " + strings.Join(e.Reasons, "; ")
}

type NotFound struct{}

func (e NotFound) Error() string {
	return "company not found"
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
