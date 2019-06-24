package administrator

import (
	"github.com/iot-my-world/brain/party/individual"
	"github.com/iot-my-world/brain/search/identifier"
	"github.com/iot-my-world/brain/security/claims"
)

type Administrator interface {
	Create(request *CreateRequest) (*CreateResponse, error)
	UpdateAllowedFields(request *UpdateAllowedFieldsRequest) (*UpdateAllowedFieldsResponse, error)
}

type CreateRequest struct {
	Claims     claims.Claims
	Individual individual.Individual
}

type CreateResponse struct {
	Individual individual.Individual
}

type UpdateAllowedFieldsRequest struct {
	Claims     claims.Claims
	Individual individual.Individual
}

type UpdateAllowedFieldsResponse struct {
	Individual individual.Individual
}

type HeartbeatRequest struct {
	Claims          claims.Claims
	SF001Identifier identifier.Identifier
}

type HeartbeatResponse struct {
}
