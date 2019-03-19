package exception

import "strings"

type UserCreation struct {
	Reasons []string
}

func (e UserCreation) Error() string {
	return "error creating user: " + strings.Join(e.Reasons, "; ")
}

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

type ChangePassword struct {
	Reasons []string
}

func (e ChangePassword) Error() string {
	return "change password error: " + strings.Join(e.Reasons, "; ")
}
