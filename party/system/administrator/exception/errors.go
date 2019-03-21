package exception

import "strings"

type SystemRetrieval struct {
	Reasons []string
}

func (e SystemRetrieval) Error() string {
	return "system retrieval error: " + strings.Join(e.Reasons, "; ")
}

type AllowedFieldsUpdate struct {
	Reasons []string
}

func (e AllowedFieldsUpdate) Error() string {
	return "allowed fields update error: " + strings.Join(e.Reasons, "; ")
}
