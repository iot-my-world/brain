package jsonRpc

import (
	jsonRpcClient "github.com/iot-my-world/brain/communication/jsonRpc/client"
	brainException "github.com/iot-my-world/brain/exception"
	"github.com/iot-my-world/brain/log"
	companyAdministrator "github.com/iot-my-world/brain/party/company/administrator"
	companyAdministratorJsonRpcAdaptor "github.com/iot-my-world/brain/party/company/administrator/adaptor/jsonRpc"
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
		"CompanyAdministrator.Create",
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
		"CompanyAdministrator.UpdateAllowedFields",
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
