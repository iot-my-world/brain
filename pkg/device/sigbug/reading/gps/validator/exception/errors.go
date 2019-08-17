package exception

import "strings"

type Validate struct {
	Reasons []string
}

func (e Validate) Error() string {
	return "error validating reading: " + strings.Join(e.Reasons, "; ")
}
