package exception

import (
	"fmt"
	"strings"
)

type RetrievingSystem struct {
	Reasons []string
}

func (e RetrievingSystem) Error() string {
	return fmt.Sprintf("error retrieving system: %s", strings.Join(e.Reasons, "; "))
}

type RetrievingCompany struct {
	Reasons []string
}

func (e RetrievingCompany) Error() string {
	return fmt.Sprintf("error retrieving company: %s", strings.Join(e.Reasons, "; "))
}

type RetrievingClient struct {
	Reasons []string
}

func (e RetrievingClient) Error() string {
	return fmt.Sprintf("error retrieving client: %s", strings.Join(e.Reasons, "; "))
}

type CollectingDevices struct {
	Reasons []string
}

func (e CollectingDevices) Error() string {
	return fmt.Sprintf("error collecting devices: %s", strings.Join(e.Reasons, "; "))
}

type CollectingReadings struct {
	Reasons []string
}

func (e CollectingReadings) Error() string {
	return fmt.Sprintf("error collecting readings: %s", strings.Join(e.Reasons, "; "))
}
