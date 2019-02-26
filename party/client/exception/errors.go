package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error() string {
	return "initial client setup error: " + strings.Join(e.Reasons, "; ")
}
