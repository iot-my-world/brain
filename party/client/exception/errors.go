package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error() string {
	return "initial client setup error: " + strings.Join(e.Reasons, "; ")
}

type NotFound struct {}

func (e NotFound) Error() string {
	return "client not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "client creation error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "client update error: " + strings.Join(e.Reasons, "; ")
}