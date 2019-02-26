package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error() string {
	return "initial role setup error: " + strings.Join(e.Reasons, "; ")
}
