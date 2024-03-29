package tracking

import (
	sigbugReading "github.com/iot-my-world/brain/pkg/device/sigbug/reading/gps"
	"github.com/iot-my-world/brain/pkg/search/identifier/party"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
)

type Report interface {
	Live(request *LiveRequest) (*LiveResponse, error)
	Historical(request *HistoricalRequest) (*HistoricalResponse, error)
}

const ServiceProvider = "Tracking-Report"
const LiveService = ServiceProvider + ".Live"
const HistoricalService = ServiceProvider + ".Historical"

var SystemUserPermissions = make([]api.Permission, 0)

var CompanyAdminUserPermissions = []api.Permission{
	LiveService,
	HistoricalService,
}

var CompanyUserPermissions = []api.Permission{
	LiveService,
	HistoricalService,
}

var ClientAdminUserPermissions = []api.Permission{
	LiveService,
	HistoricalService,
}

var ClientUserPermissions = []api.Permission{
	LiveService,
	HistoricalService,
}

type LiveRequest struct {
	Claims           claims.Claims
	PartyIdentifiers []party.Identifier
}

type LiveResponse struct {
	ZX303TrackerGPSReadings []sigbugReading.Reading
}

type HistoricalRequest struct {
	Claims claims.Claims
}

type HistoricalResponse struct {
	ZX303TrackerGPSReadings []sigbugReading.Reading
}
