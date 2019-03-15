package exception

import "strings"

type InvalidClaims struct {
	Reasons []string
}

func (e InvalidClaims) Error() string {
	return "invalid claims: " + strings.Join(e.Reasons, "; ")
}

type UserRetrieval struct {
	Reasons []string
}

func (e UserRetrieval) Error() string {
	return "user retrieval error: " + strings.Join(e.Reasons, "; ")
}
