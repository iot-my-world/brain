package exception

import "strings"

type InvalidParty struct {
	Reasons []string
}

func (e InvalidParty) Error() string {
	return "invalid party type: " + strings.Join(e.Reasons, "; ")
}

type PartyRetrieval struct {
	Reasons []string
}

func (e PartyRetrieval) Error() string {
	return "party retrieval error: " + strings.Join(e.Reasons, "; ")
}

type InvalidClaims struct {
	Reasons []string
}

func (e InvalidClaims) Error() string {
	return "invalid claims: " + strings.Join(e.Reasons, "; ")
}
