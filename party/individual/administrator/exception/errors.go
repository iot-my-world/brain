package exception

import (
	"strings"
)

type IndividualRetrieval struct {
	Reasons []string
}

func (e IndividualRetrieval) Error() string {
	return "error retrieving individual: " + strings.Join(e.Reasons, "; ")
}

type ReadingCollection struct {
	Reasons []string
}

func (e ReadingCollection) Error() string {
	return "error collecting readings : " + strings.Join(e.Reasons, "; ")
}

type IndividualUpdate struct {
	Reasons []string
}

func (e IndividualUpdate) Error() string {
	return "error updating individual: " + strings.Join(e.Reasons, "; ")
}

type ReadingUpdate struct {
	Reasons []string
}

func (e ReadingUpdate) Error() string {
	return "error updating reading: " + strings.Join(e.Reasons, "; ")
}

type IndividualCreation struct {
	Reasons []string
}

func (e IndividualCreation) Error() string {
	return "error creating individual: " + strings.Join(e.Reasons, "; ")
}

type Heartbeat struct {
	Reasons []string
}

func (e Heartbeat) Error() string {
	return "heartbeat error: " + strings.Join(e.Reasons, "; ")
}
