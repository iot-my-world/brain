package exception

import "strings"

type NotFound struct{}

func (e NotFound) Error() string {
	return "party not found"
}

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

type CreateAndInviteCompany struct {
	Reasons []string
}

func (e CreateAndInviteCompany) Error() string {
	return "create and invite company error: " + strings.Join(e.Reasons, "; ")
}

type CreateAndInviteClient struct {
	Reasons []string
}

func (e CreateAndInviteClient) Error() string {
	return "create and invite client error: " + strings.Join(e.Reasons, "; ")
}
