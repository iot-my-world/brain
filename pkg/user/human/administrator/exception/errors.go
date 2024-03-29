package exception

import "strings"

type UserCreation struct {
	Reasons []string
}

func (e UserCreation) Error() string {
	return "error creating user: " + strings.Join(e.Reasons, "; ")
}

type InvalidClaims struct {
	Reasons []string
}

func (e InvalidClaims) Error() string {
	return "invalid claims: " + strings.Join(e.Reasons, "; ")
}

type UserRetrieval struct {
	Reasons []string
}

func (e UserRetrieval) Error() string {
	return "user retrieval error: " + strings.Join(e.Reasons, "; ")
}

type SetPassword struct {
	Reasons []string
}

func (e SetPassword) Error() string {
	return "set password error: " + strings.Join(e.Reasons, "; ")
}

type AllowedFieldsUpdate struct {
	Reasons []string
}

func (e AllowedFieldsUpdate) Error() string {
	return "allowed fields update error: " + strings.Join(e.Reasons, "; ")
}

type CheckPassword struct {
	Reasons []string
}

func (e CheckPassword) Error() string {
	return "check password error: " + strings.Join(e.Reasons, "; ")
}

type TokenGeneration struct {
	Reasons []string
}

func (e TokenGeneration) Error() string {
	return "token generation error: " + strings.Join(e.Reasons, "; ")
}

type EmailGeneration struct {
	Reasons []string
}

func (e EmailGeneration) Error() string {
	return "email generation error: " + strings.Join(e.Reasons, "; ")
}

type UpdateAllowedFields struct {
	Reasons []string
}

func (e UpdateAllowedFields) Error() string {
	return "error updating allowed fields: " + strings.Join(e.Reasons, "; ")
}

type GetMyUser struct {
	Reasons []string
}

func (e GetMyUser) Error() string {
	return "error getting my user: " + strings.Join(e.Reasons, "; ")
}

type Create struct {
	Reasons []string
}

func (e Create) Error() string {
	return "error creating user: " + strings.Join(e.Reasons, "; ")
}

type UpdatePassword struct {
	Reasons []string
}

func (e UpdatePassword) Error() string {
	return "error updating password: " + strings.Join(e.Reasons, "; ")
}

type ForgotPassword struct {
	Reasons []string
}

func (e ForgotPassword) Error() string {
	return "error processing forgotten password: " + strings.Join(e.Reasons, "; ")
}
