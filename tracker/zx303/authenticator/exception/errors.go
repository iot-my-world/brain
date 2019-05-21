package exception

import "strings"

type Login struct {
	Reasons []string
}

func (e Login) Error() string {
	return "login error: " + strings.Join(e.Reasons, "; ")
}

type Logout struct {
	Reasons []string
}

func (e Logout) Error() string {
	return "Logout error: " + strings.Join(e.Reasons, "; ")
}
