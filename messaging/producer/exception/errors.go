package exception

import "strings"

type Start struct {
	Reasons []string
}

func (e Start) Error() string {
	return "starting error: " + strings.Join(e.Reasons, "; ")
}

type Produce struct {
	Reasons []string
}

func (e Produce) Error() string {
	return "production error: " + strings.Join(e.Reasons, "; ")
}
