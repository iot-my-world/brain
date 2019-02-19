package exception

import (
	"fmt"
	"strings"
)

type InvalidJWT struct {
	Reasons []string
}

func (e InvalidJWT) Error() string {
	return fmt.Sprintf("invalid JWT: %s", strings.Join(e.Reasons, "; "))
}

type JWTVerification struct {
	Reasons []string
}

func (e JWTVerification) Error() string {
	return fmt.Sprintf("JWT verification error: %s", strings.Join(e.Reasons, ";"))
}

type JWTUnmarshalling struct {
	Reasons []string
}

func (e JWTUnmarshalling) Error() string {
	return fmt.Sprintf("JWT unmarshalling error: %s", strings.Join(e.Reasons, "; "))
}
