package exception

import (
	"strings"
)

type ZX303TaskSubmission struct {
	Reasons []string
}

func (e ZX303TaskSubmission) Error() string {
	return "error submitting zx303 task: " + strings.Join(e.Reasons, "; ")
}
