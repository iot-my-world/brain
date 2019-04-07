package jsonRpc

import (
	deviceAdministrator "gitlab.com/iotTracker/brain/tracker/device/administrator"
)

type adaptor struct {
	deviceAdministrator deviceAdministrator.Administrator
}

func New(
	deviceAdministrator deviceAdministrator.Administrator,
) *adaptor {
	return &adaptor{
		deviceAdministrator: deviceAdministrator,
	}
}
