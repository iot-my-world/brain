package exception

import (
	"fmt"
	"strings"
)

type GroupCreation struct {
	GroupName string
	Reasons   []string
}

func (e GroupCreation) Error() string {
	return fmt.Sprintf("error creating consumer group %s: %s", e.GroupName, strings.Join(e.Reasons, "; "))
}

type Consumption struct {
	Reasons []string
}

func (e Consumption) Error() string {
	return "error consuming from kafka: " + strings.Join(e.Reasons, "; ")
}
