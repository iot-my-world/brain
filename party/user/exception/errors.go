package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error() string {
	return "initial user setup error: " + strings.Join(e.Reasons, "; ")
}

type NotFound struct {}

func (e NotFound) Error() string {
	return "user not found"
}
