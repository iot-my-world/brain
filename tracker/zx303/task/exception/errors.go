package exception

import "strings"

type ExecutingStep struct {
	Reasons []string
}

func (e ExecutingStep) Error() string {
	return "error getting executing step: " + strings.Join(e.Reasons, "; ")
}

type NextStep struct {
	Reasons []string
}

func (e NextStep) Error() string {
	return "error getting next step: " + strings.Join(e.Reasons, "; ")
}
