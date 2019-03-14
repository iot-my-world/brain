package exception

import (
	"strings"
)

// DeviceRetrieval error
type DeviceRetrieval struct {
	Reasons []string
}

func (e DeviceRetrieval) Error() string {
	return "error retrieving device: " + strings.Join(e.Reasons, "; ")
}

// ReadingCollection error
type ReadingCollection struct {
	Reasons []string
}

func (e ReadingCollection) Error() string {
	return "error collecting readings : " + strings.Join(e.Reasons, "; ")
}

// DeviceUpdate error
type DeviceUpdate struct {
	Reasons []string
}

func (e DeviceUpdate) Error() string {
	return "error updating device: " + strings.Join(e.Reasons, "; ")
}

// ReadingUpdate error
type ReadingUpdate struct {
	Reasons []string
}

func (e ReadingUpdate) Error() string {
	return "error updating reading: " + strings.Join(e.Reasons, "; ")
}
