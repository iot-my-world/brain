package exception

import "strings"

type HandleDataMessage struct {
	Reasons []string
}

func (e HandleDataMessage) Error() string {
	return "error handling data message: " + strings.Join(e.Reasons, "; ")
}
