package exception

import (
	"strings"
)

type ZX303GPSReadingCreation struct {
	Reasons []string
}

func (e ZX303GPSReadingCreation) Error() string {
	return "error creating zx303 gps reading: " + strings.Join(e.Reasons, "; ")
}
