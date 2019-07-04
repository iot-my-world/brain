package exception

import "strings"

type Wrapping struct {
	Reasons []string
}

func (e Wrapping) Error() string {
	return "error wrapping party: " + strings.Join(e.Reasons, "; ")
}

type InvalidPartyType struct {
	Reasons []string
}

func (e InvalidPartyType) Error() string {
	return "invalid party type: " + strings.Join(e.Reasons, "; ")
}

type Unwrapping struct {
	Reasons []string
}

func (e Unwrapping) Error() string {
	return "party unwrapping error: " + strings.Join(e.Reasons, "; ")
}
