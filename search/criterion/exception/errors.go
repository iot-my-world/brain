package exception

import (
	"fmt"
	"strings"
)

type Unwrapping struct {
	Reasons []string
}

func (e Unwrapping) Error() string {
	return fmt.Sprintf("unwrapping criterion: %s", strings.Join(e.Reasons, "; "))
}

type Invalid struct {
	Reasons []string
}

func (e Invalid) Error() string {
	return fmt.Sprintf("invalid critrion: %s", strings.Join(e.Reasons, "; "))
}
