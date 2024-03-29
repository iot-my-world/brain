package exception

import (
	"strings"
)

type DeviceRetrieval struct {
	Reasons []string
}

func (e DeviceRetrieval) Error() string {
	return "error retrieving message: " + strings.Join(e.Reasons, "; ")
}

type ReadingCollection struct {
	Reasons []string
}

func (e ReadingCollection) Error() string {
	return "error collecting readings : " + strings.Join(e.Reasons, "; ")
}

type DeviceUpdate struct {
	Reasons []string
}

func (e DeviceUpdate) Error() string {
	return "error updating message: " + strings.Join(e.Reasons, "; ")
}

type ReadingUpdate struct {
	Reasons []string
}

func (e ReadingUpdate) Error() string {
	return "error updating reading: " + strings.Join(e.Reasons, "; ")
}

type DeviceCreation struct {
	Reasons []string
}

func (e DeviceCreation) Error() string {
	return "error creating message: " + strings.Join(e.Reasons, "; ")
}

type Heartbeat struct {
	Reasons []string
}

func (e Heartbeat) Error() string {
	return "heartbeat error: " + strings.Join(e.Reasons, "; ")
}
