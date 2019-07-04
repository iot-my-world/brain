package exception

import "strings"

type CompanyCreation struct {
	Reasons []string
}

func (e CompanyCreation) Error() string {
	return "company creation error: " + strings.Join(e.Reasons, "; ")
}

type CompanyRetrieval struct {
	Reasons []string
}

func (e CompanyRetrieval) Error() string {
	return "company retrieval error: " + strings.Join(e.Reasons, "; ")
}

type AllowedFieldsUpdate struct {
	Reasons []string
}

func (e AllowedFieldsUpdate) Error() string {
	return "allowed fields update error: " + strings.Join(e.Reasons, "; ")
}

type Delete struct {
	Reasons []string
}

func (e Delete) Error() string {
	return "delete company error: " + strings.Join(e.Reasons, "; ")
}
