package exception

import "strings"

type InvalidData struct {
	Reasons []string
}

func (e InvalidData) Error() string {
	return "invalid email data: " + strings.Join(e.Reasons, "; ")
}

type TemplateExecution struct {
	Reasons []string
}

func (e TemplateExecution) Error() string {
	return "error executing email template: " + strings.Join(e.Reasons, "; ")
}
