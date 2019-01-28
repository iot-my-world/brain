package exception

import "strings"

type RequestInvalid struct {
	Reasons []string
}


func (e RequestInvalid) Error() string {
	return "invalidRequest: " + strings.Join(e.Reasons, ";")
}