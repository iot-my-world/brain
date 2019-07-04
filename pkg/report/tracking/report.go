package tracking

import (
	"github.com/iot-my-world/brain/pkg/search/identifier/party"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/permission/api"
	zx303TrackerGPSReading "github.com/iot-my-world/brain/pkg/tracker/zx303/reading/gps"
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
	ZX303TrackerGPSReadings []zx303TrackerGPSReading.Reading
}

type HistoricalRequest struct {
	Claims claims.Claims
}

type HistoricalResponse struct {
	ZX303TrackerGPSReadings []zx303TrackerGPSReading.Reading
}
