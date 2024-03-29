package exception

import "strings"

type RequestInvalid struct {
	Reasons []string
}

func (e RequestInvalid) Error() string {
	return "invalidRequest: " + strings.Join(e.Reasons, "; ")
}

type Unexpected struct {
	Reasons []string
}

func (e Unexpected) Error() string {
	return "unexpected error: " + strings.Join(e.Reasons, "; ")
}

type UUIDGeneration struct {
	Reasons []string
}

func (e UUIDGeneration) Error() string {
	return "UUID generation error: " + strings.Join(e.Reasons, "; ")
}

type NotImplemented struct {
}

func (e NotImplemented) Error() string {
	return "not implemented"
}
