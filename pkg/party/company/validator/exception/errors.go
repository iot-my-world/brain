package exception

import "strings"

type Validate struct {
	Reasons []string
}

func (e Validate) Error() string {
	return "error validating company: " + strings.Join(e.Reasons, "; ")
}
