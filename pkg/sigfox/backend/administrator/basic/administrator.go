package basic

import (
	"crypto/rsa"
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	sigfoxBackendClaims "github.com/iot-my-world/brain/pkg/security/claims/sigfoxBackend"
	"github.com/iot-my-world/brain/pkg/security/token"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/action"
	backendAdministrator "github.com/iot-my-world/brain/pkg/sigfox/backend/administrator"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/administrator/exception"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/recordHandler"
	"github.com/iot-my-world/brain/pkg/sigfox/backend/validator"
)

type administrator struct {
	backendDeviceValidator validator.Validator
	backendRecordHandler   recordHandler.RecordHandler
	jwtGenerator           token.JWTGenerator
}

func New(
	backendDeviceValidator validator.Validator,
	backendRecordHandler recordHandler.RecordHandler,
	rsaPrivateKey *rsa.PrivateKey,
) backendAdministrator.Administrator {
	return &administrator{
		backendDeviceValidator: backendDeviceValidator,
		backendRecordHandler:   backendRecordHandler,
		jwtGenerator:           token.NewJWTGenerator(rsaPrivateKey),
	}
}

func (a *administrator) ValidateCreateRequest(request *backendAdministrator.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		backendDeviceValidateResponse, err := a.backendDeviceValidator.Validate(&validator.ValidateRequest{
			Claims:  request.Claims,
			Backend: request.Backend,
			Action:  action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating backend backend: "+err.Error())
		} else {
			if len(backendDeviceValidateResponse.ReasonsInvalid) > 0 {
				for _, reason := range backendDeviceValidateResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("backend backend invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (a *administrator) Create(request *backendAdministrator.CreateRequest) (*backendAdministrator.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	createResponse, err := a.backendRecordHandler.Create(&recordHandler.CreateRequest{
		Backend: request.Backend,
	})
	if err != nil {
		err = exception.DeviceCreation{Reasons: []string{"creation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	claimsForToken := sigfoxBackendClaims.SigfoxBackend{
		BackendId: id.Identifier{Id: createResponse.Backend.Id},
	}
	backendToken, err := a.jwtGenerator.GenerateToken(claimsForToken)
	if err != nil {
		err = exception.DeviceCreation{Reasons: []string{"token generation", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	createResponse.Backend.Token = backendToken
	if _, err := a.backendRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: createResponse.Backend.Id},
		Backend:    createResponse.Backend,
	}); err != nil {
		err = exception.DeviceCreation{Reasons: []string{"update sigfox backend with token", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &backendAdministrator.CreateResponse{
		Backend: createResponse.Backend,
	}, nil
}

func (a *administrator) ValidateUpdateAllowedFieldsRequest(request *backendAdministrator.UpdateAllowedFieldsRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	} else {
		// backend must be valid
		validationResponse, err := a.backendDeviceValidator.Validate(&validator.ValidateRequest{
			Claims: request.Claims,
			Action: action.UpdateAllowedFields,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating backend: "+err.Error())
		} else {
			if len(validationResponse.ReasonsInvalid) > 0 {
				for _, reason := range validationResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("backend backend invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) UpdateAllowedFields(request *backendAdministrator.UpdateAllowedFieldsRequest) (*backendAdministrator.UpdateAllowedFieldsResponse, error) {
	if err := a.ValidateUpdateAllowedFieldsRequest(request); err != nil {
		return nil, err
	}

	// retrieve the backend
	backendRetrieveResponse, err := a.backendRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Backend.Id},
	})
	if err != nil {
		return nil, exception.DeviceRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the backend

	// update the backend
	_, err = a.backendRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Backend.Id},
		Backend:    backendRetrieveResponse.Backend,
	})
	if err != nil {
		return nil, exception.DeviceUpdate{Reasons: []string{err.Error()}}
	}

	return &backendAdministrator.UpdateAllowedFieldsResponse{
		Backend: backendRetrieveResponse.Backend,
	}, nil
}
