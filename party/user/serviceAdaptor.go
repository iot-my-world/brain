package user

import (
	"net/http"
)

type serviceAdaptor struct{
	RecordHandler
}

func NewServiceAdaptor(recordHandler RecordHandler) *serviceAdaptor {
	return &serviceAdaptor{
		recordHandler,
	}
}

type CreateUserRequest struct {
	NewUser NewUser `json:"newUser"`
}

func (s *serviceAdaptor) Create(r * http.Request, request *CreateUserRequest, response *CreateResponse) error {
	createRequest := &CreateRequest{
		NewUser{
			// Personal Details
			Name: request.NewUser.Name,
			Surname: request.NewUser.Surname,
			IDNo: request.NewUser.IDNo,

			// System Details
			Username: request.NewUser.Username,
			Password: request.NewUser.Password,
			SystemRole: request.NewUser.SystemRole,
		},

	}
	return s.RecordHandler.Create(createRequest, response)
}

func (s *serviceAdaptor) RetrieveAll(r *http.Request, request *RetrieveAllRequest, response *RetrieveAllResponse) error {

	return s.RecordHandler.RetrieveAll(request, response)
}

func (s *serviceAdaptor) Update(r *http.Request, request *UpdateRequest, response *UpdateResponse) error {
	return s.RecordHandler.Update(request, response)
}

func (s *serviceAdaptor) Delete(r *http.Request, request *DeleteRequest, response *DeleteResponse) error {
	return s.RecordHandler.Delete(request, response)
}