package exception

import "strings"

type Retrieval struct {
	Reasons []string
}

func (e Retrieval) Error() string {
	return "retrieval error: " + strings.Join(e.Reasons, "; ")
}
