package exception

import "strings"

type NotEnoughRowsInSheet struct {
	Reasons []string
}

func (e NotEnoughRowsInSheet) Error() string {
	return "not enough rows in sheet: " + strings.Join(e.Reasons, "; ")
}

type OpeningFile struct {
	Reasons []string
}

func (e OpeningFile) Error() string {
	return "error opening file: " + strings.Join(e.Reasons, "; ")
}
