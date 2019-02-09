package exception

import (
	"fmt"
	"strings"
)

type Unwrapping struct {
	Reasons []string
}

func (e Unwrapping) Error() string {
	return fmt.Sprintf("unwrapping identifier: %s", strings.Join(e.Reasons, "; "))
}

type InvalidIdentifier struct {
	Reasons []string
}

func (e InvalidIdentifier) Error() string {
	return fmt.Sprintf("invalid identifier: %s", strings.Join(e.Reasons, "; "))
}
