package exception

import "strings"

type Wrapping struct {
	Reasons []string
}

func (e Wrapping) Error() string {
	return "device wrapping error: " + strings.Join(e.Reasons, "; ")
}

type UnWrapping struct {
	Reasons []string
}

func (e UnWrapping) Error() string {
	return "device unwrapping error: " + strings.Join(e.Reasons, "; ")
}
