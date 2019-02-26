package exception

import "strings"

type NotFound struct{}

func (e NotFound) Error() string {
	return "user not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "user creation error: " + strings.Join(e.Reasons, "; ")
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "user update error: " + strings.Join(e.Reasons, "; ")
}

type ChangePassword struct {
	Reasons []string
}

func (e ChangePassword) Error() string {
	return "change password error: " + strings.Join(e.Reasons, "; ")
}
