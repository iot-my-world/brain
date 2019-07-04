package exception

import (
	"strings"
)

type APIUserRetrieval struct {
	Reasons []string
}

func (e APIUserRetrieval) Error() string {
	return "error retrieving api user: " + strings.Join(e.Reasons, "; ")
}

type ReadingCollection struct {
	Reasons []string
}

func (e ReadingCollection) Error() string {
	return "error collecting readings : " + strings.Join(e.Reasons, "; ")
}

type APIUserUpdate struct {
	Reasons []string
}

func (e APIUserUpdate) Error() string {
	return "error updating api user: " + strings.Join(e.Reasons, "; ")
}

type ReadingUpdate struct {
	Reasons []string
}

func (e ReadingUpdate) Error() string {
	return "error updating reading: " + strings.Join(e.Reasons, "; ")
}

type APIUserCreation struct {
	Reasons []string
}

func (e APIUserCreation) Error() string {
	return "error creating api user: " + strings.Join(e.Reasons, "; ")
}

type PasswordGeneration struct {
	Reasons []string
}

func (e PasswordGeneration) Error() string {
	return "error generating api user password: " + strings.Join(e.Reasons, "; ")
}

type PasswordHash struct {
	Reasons []string
}

func (e PasswordHash) Error() string {
	return "error hashing api user password: " + strings.Join(e.Reasons, "; ")
}
