package exception

import (
	"fmt"
	"strings"
)

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

type SheetDoesNotExist struct {
	SheetName string
}

func (e SheetDoesNotExist) Error() string {
	return fmt.Sprintf("sheet with name %s does not exist", e.SheetName)
}
