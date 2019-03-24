package exception

import "strings"

type NotFound struct{}

func (e NotFound) Error() string {
	return "reading not found"
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "reading creation error: " + strings.Join(e.Reasons, "; ")
}

type CreateBulk struct {
	Reason string
}

func (e CreateBulk) Error() string {
	return "bulk creation error: " + e.Reason
}

type Update struct {
	Reasons []string
}

func (e Update) Error() string {
	return "reading update error: " + strings.Join(e.Reasons, "; ")
}
