package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	clientAdministrator "gitlab.com/iotTracker/brain/party/client/administrator"
	clientAdministratorException "gitlab.com/iotTracker/brain/party/client/administrator/exception"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	"gitlab.com/iotTracker/brain/search/identifier/id"
)

type basicAdministrator struct {
	clientRecordHandler clientRecordHandler.RecordHandler
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
) clientAdministrator.Administrator {
	return &basicAdministrator{
		clientRecordHandler: clientRecordHandler,
	}
}

func (ba *basicAdministrator) ValidateUpdateAllowedFieldsRequest(request *clientAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (ba *basicAdministrator) UpdateAllowedFields(request *clientAdministrator.UpdateAllowedFieldsRequest, response *clientAdministrator.UpdateAllowedFieldsResponse) error {
	if err := ba.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return err
	}

	// retrieve the client
	clientRetrieveResponse := clientRecordHandler.RetrieveResponse{}
	if err := ba.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
	}, &clientRetrieveResponse); err != nil {
		return clientAdministratorException.ClientRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the client
	//clientRetrieveResponse.Client.Id = request.Client.Id
	//clientRetrieveResponse.Client.ParentId = request.Client.ParentId
	//clientRetrieveResponse.Client.ParentPartyType = request.Client.ParentPartyType
	clientRetrieveResponse.Client.Name = request.Client.Name
	//clientRetrieveResponse.Client.AdminEmailAddress = request.Client.AdminEmailAddress

	// update the client
	clientUpdateResponse := clientRecordHandler.UpdateResponse{}
	if err := ba.clientRecordHandler.Update(&clientRecordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Client.Id},
		Client:     clientRetrieveResponse.Client,
	}, &clientUpdateResponse); err != nil {
		return clientAdministratorException.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	response.Client = clientUpdateResponse.Client

	return nil
}
