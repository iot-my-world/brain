package exception

import "strings"

type HandleGPSMessage struct {
	Reasons []string
}

func (e HandleGPSMessage) Error() string {
	return "error handling gps message: " + strings.Join(e.Reasons, "; ")
}
