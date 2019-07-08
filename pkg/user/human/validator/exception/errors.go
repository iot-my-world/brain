package exception

import "strings"

type Validate struct {
	Reasons []string
}

func (e Validate) Error() string {
	return "user validation error: " + strings.Join(e.Reasons, "; ")
}
