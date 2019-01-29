package exception

import "strings"

type InitialSetup struct {
	Reasons []string
}

func (e InitialSetup) Error () string {
	return "initial user setup error: " + strings.Join()
}
