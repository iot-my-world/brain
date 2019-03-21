package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	systemAdministrator "gitlab.com/iotTracker/brain/party/system/administrator"
	systemAdministratorException "gitlab.com/iotTracker/brain/party/system/administrator/exception"
	systemRecordHandler "gitlab.com/iotTracker/brain/party/system/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type administrator struct {
	systemRecordHandler systemRecordHandler.RecordHandler
}

func New(
	systemRecordHandler systemRecordHandler.RecordHandler,
) systemAdministrator.Administrator {
	return &administrator{
		systemRecordHandler: systemRecordHandler,
	}
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *systemAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *systemAdministrator.UpdateAllowedFieldsRequest) (*systemAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the system
	systemRetrieveResponse, err := a.systemRecordHandler.Retrieve(&systemRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.System.Id},
	})
	if err != nil {
		return nil, systemAdministratorException.SystemRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the system
	//systemRetrieveResponse.System.Id = request.System.Id
	//systemRetrieveResponse.System.ParentId = request.System.ParentId
	//systemRetrieveResponse.System.ParentPartyType = request.System.ParentPartyType
	//systemRetrieveResponse.System.ParentId = request.System.ParentId
	systemRetrieveResponse.System.Name = request.System.Name
	//systemRetrieveResponse.System.AdminEmailAddress = request.System.AdminEmailAddress

	// update the system
	systemUpdateResponse, err := a.systemRecordHandler.Update(&systemRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.System.Id},
		System:     systemRetrieveResponse.System,
	})
	if err != nil {
		return systemAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	response.System = systemUpdateResponse.System

	return nil
}
