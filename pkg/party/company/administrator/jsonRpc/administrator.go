package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/internal/api/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	companyAdministrator "github.com/iot-my-world/brain/pkg/party/company/administrator"
	companyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/pkg/party/company/administrator/adaptor/jsonRpc"
	wrappedIdentifier "github.com/iot-my-world/brain/pkg/search/identifier/wrapped"
)

type administrator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) companyAdministrator.Administrator {
	return &administrator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *administrator) ValidateCreateRequest(request *companyAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *companyAdministrator.CreateRequest) (*companyAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	companyCreateResponse := companyAdministratorJsonRpcAdaptor.CreateResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		companyAdministrator.CreateService,
		companyAdministratorJsonRpcAdaptor.CreateRequest{
			Company: request.Company,
		},
		&companyCreateResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &companyAdministrator.CreateResponse{Company: companyCreateResponse.Company}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *companyAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *companyAdministrator.UpdateAllowedFieldsRequest) (*companyAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	companyUpdateAllowedFieldsResponse := companyAdministratorJsonRpcAdaptor.UpdateAllowedFieldsResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		companyAdministrator.UpdateAllowedFieldsService,
		companyAdministratorJsonRpcAdaptor.UpdateAllowedFieldsRequest{
			Company: request.Company,
		},
		&companyUpdateAllowedFieldsResponse,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &companyAdministrator.UpdateAllowedFieldsResponse{
		Company: companyUpdateAllowedFieldsResponse.Company,
	}, nil
}

func (a *administrator) ValidateDeleteRequest(request *companyAdministrator.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.CompanyIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "company identifier is nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Delete(request *companyAdministrator.DeleteRequest) (*companyAdministrator.DeleteResponse, error) {
	if err := a.ValidateDeleteRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// wrap identifier
	id, err := wrappedIdentifier.Wrap(request.CompanyIdentifier)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	response := companyAdministratorJsonRpcAdaptor.DeleteResponse{}
	if err := a.jsonRpcClient.JsonRpcRequest(
		companyAdministrator.DeleteService,
		companyAdministratorJsonRpcAdaptor.DeleteRequest{
			CompanyIdentifier: *id,
		},
		&response); err != nil {
		return nil, err
	}

	return &companyAdministrator.DeleteResponse{}, nil
}
