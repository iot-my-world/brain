package exception

import (
	"strings"
)

type ZX303StatusReadingCreation struct {
	Reasons []string
}

func (e ZX303StatusReadingCreation) Error() string {
	return "error creating zx303 status reading: " + strings.Join(e.Reasons, "; ")
}
