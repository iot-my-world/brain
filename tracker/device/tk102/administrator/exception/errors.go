package exception

import (
	"strings"
)

type DeviceRetrieval struct {
	Reasons []string
}

func (e DeviceRetrieval) Error() string {
	return "error retrieving device: " + strings.Join(e.Reasons, "; ")
}
