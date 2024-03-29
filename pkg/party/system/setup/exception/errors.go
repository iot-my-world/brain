package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error() string {
	return "initial system setup error: " + strings.Join(e.Reasons, "; ")
}
