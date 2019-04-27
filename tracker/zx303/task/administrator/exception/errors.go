package exception

import (
	"strings"
)

type ZX303TaskCreation struct {
	Reasons []string
}

func (e ZX303TaskCreation) Error() string {
	return "error creating zx303 task: " + strings.Join(e.Reasons, "; ")
}
