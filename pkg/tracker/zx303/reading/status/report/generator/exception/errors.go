package exception

import "strings"

type BatteryReportGeneration struct {
	Reasons []string
}

func (e BatteryReportGeneration) Error() string {
	return "error generating battery report: " + strings.Join(e.Reasons, "; ")
}
