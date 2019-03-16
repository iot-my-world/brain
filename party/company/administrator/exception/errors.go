package exception

import "strings"

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
