package basic

import (
	brainException "github.com/iot-my-world/brain/exception"
	administrator2 "github.com/iot-my-world/brain/pkg/party/system/administrator"
	"github.com/iot-my-world/brain/pkg/party/system/administrator/exception"
	"github.com/iot-my-world/brain/pkg/party/system/recordHandler"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
)

type administrator struct {
	systemRecordHandler recordHandler.RecordHandler
}

func New(
	systemRecordHandler recordHandler.RecordHandler,
) administrator2.Administrator {
	return &administrator{
		systemRecordHandler: systemRecordHandler,
	}
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *administrator2.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *administrator2.UpdateAllowedFieldsRequest) (*administrator2.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the system
	systemRetrieveResponse, err := a.systemRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.System.Id},
	})
	if err != nil {
		return nil, exception.SystemRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the system
	//systemRetrieveResponse.System.Id = request.System.Id
	//systemRetrieveResponse.System.ParentId = request.System.ParentId
	//systemRetrieveResponse.System.ParentPartyType = request.System.ParentPartyType
	//systemRetrieveResponse.System.ParentId = request.System.ParentId
	systemRetrieveResponse.System.Name = request.System.Name
	//systemRetrieveResponse.System.AdminEmailAddress = request.System.AdminEmailAddress

	// update the system
	systemUpdateResponse, err := a.systemRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.System.Id},
		System:     systemRetrieveResponse.System,
	})
	if err != nil {
		return exception.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	response.System = systemUpdateResponse.System

	return nil
}
