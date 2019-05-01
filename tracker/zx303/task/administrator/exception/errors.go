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

type ZX303TaskFail struct {
	Reasons []string
}

func (e ZX303TaskFail) Error() string {
	return "failed to transition ZX303 task to fail: " + strings.Join(e.Reasons, "; ")
}
