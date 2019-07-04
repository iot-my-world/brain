package basic

import (
	"fmt"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/pkg/party"
	"github.com/iot-my-world/brain/pkg/party/company/action"
	administrator2 "github.com/iot-my-world/brain/pkg/party/company/administrator"
	"github.com/iot-my-world/brain/pkg/party/company/administrator/exception"
	"github.com/iot-my-world/brain/pkg/party/company/recordHandler"
	"github.com/iot-my-world/brain/pkg/party/company/validator"
	"github.com/iot-my-world/brain/pkg/search/criterion"
	exactTextCriterion "github.com/iot-my-world/brain/pkg/search/criterion/exact/text"
	"github.com/iot-my-world/brain/pkg/search/identifier/id"
	humanUser "github.com/iot-my-world/brain/pkg/user/human"
	userRecordHandler "github.com/iot-my-world/brain/pkg/user/human/recordHandler"
	humanUserLoginClaims "github.com/iot-my-world/brain/security/claims/login/user/human"
)

type administrator struct {
	companyRecordHandler recordHandler.RecordHandler
	companyValidator     validator.Validator
	userRecordHandler    userRecordHandler.RecordHandler
	systemClaims         *humanUserLoginClaims.Login
}

func New(
	companyRecordHandler recordHandler.RecordHandler,
	companyValidator validator.Validator,
	userRecordHandler userRecordHandler.RecordHandler,
	systemClaims *humanUserLoginClaims.Login,
) administrator2.Administrator {
	return &administrator{
		companyRecordHandler: companyRecordHandler,
		companyValidator:     companyValidator,
		userRecordHandler:    userRecordHandler,
		systemClaims:         systemClaims,
	}
}

func (a *administrator) ValidateCreateRequest(request *administrator2.CreateRequest) error {
	reasonsInvalid := make([]string, 0)

	// A new company can only be made by root
	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "nil claims")
	} else {
		if request.Claims.PartyDetails().PartyType != party.System {
			reasonsInvalid = append(reasonsInvalid, "only system party can make a new company")
		}

		// company must be valid
		validationResponse, err := a.companyValidator.Validate(&validator.ValidateRequest{
			Claims:  request.Claims,
			Company: request.Company,
			Action:  action.Create,
		})
		if err != nil {
			reasonsInvalid = append(reasonsInvalid, "error validating company: "+err.Error())
		} else {
			if len(validationResponse.ReasonsInvalid) > 0 {
				for _, reason := range validationResponse.ReasonsInvalid {
					reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("company invalid: %s - %s - %s", reason.Field, reason.Type, reason.Help))
				}
			}
		}
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Create(request *administrator2.CreateRequest) (*administrator2.CreateResponse, error) {
	if err := a.ValidateCreateRequest(request); err != nil {
		return nil, err
	}

	// create the company
	companyCreateResponse, err := a.companyRecordHandler.Create(&recordHandler.CreateRequest{
		Company: request.Company,
	})
	if err != nil {
		return nil, exception.CompanyCreation{Reasons: []string{"creating company", err.Error()}}
	}

	// create minimal admin user for the company
	if _, err := a.userRecordHandler.Create(&userRecordHandler.CreateRequest{
		User: humanUser.User{
			EmailAddress:    companyCreateResponse.Company.AdminEmailAddress,
			ParentPartyType: companyCreateResponse.Company.ParentPartyType,
			ParentId:        companyCreateResponse.Company.ParentId,
			PartyType:       party.Company,
			PartyId:         id.Identifier{Id: companyCreateResponse.Company.Id},
		},
	}); err != nil {
		return nil, exception.CompanyCreation{Reasons: []string{"creating admin user", err.Error()}}
	}

	return &administrator2.CreateResponse{Company: companyCreateResponse.Company}, nil
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

	// retrieve the company
	companyRetrieveResponse, err := a.companyRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Company.Id},
	})
	if err != nil {
		return nil, exception.CompanyRetrieval{Reasons: []string{err.Error()}}
	}

	// update the allowed fields on the company
	//companyRetrieveResponse.Company.Id = request.Company.Id
	//companyRetrieveResponse.Company.ParentId = request.Company.ParentId
	//companyRetrieveResponse.Company.ParentPartyType = request.Company.ParentPartyType
	companyRetrieveResponse.Company.Name = request.Company.Name
	//companyRetrieveResponse.Company.AdminEmailAddress = request.Company.AdminEmailAddress

	// update the company
	if _, err := a.companyRecordHandler.Update(&recordHandler.UpdateRequest{
		Claims:     request.Claims,
		Identifier: id.Identifier{Id: request.Company.Id},
		Company:    companyRetrieveResponse.Company,
	}); err != nil {
		return nil, exception.AllowedFieldsUpdate{Reasons: []string{"updating", err.Error()}}
	}

	return &administrator2.UpdateAllowedFieldsResponse{Company: companyRetrieveResponse.Company}, nil
}

func (a *administrator) ValidateDeleteRequest(request *administrator2.DeleteRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.CompanyIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "company identifier is nil")
	}

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (a *administrator) Delete(request *administrator2.DeleteRequest) (*administrator2.DeleteResponse, error) {
	if err := a.ValidateDeleteRequest(request); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// retrieve the company to be deleted
	companyRetrieveResponse, err := a.companyRecordHandler.Retrieve(&recordHandler.RetrieveRequest{
		Claims:     request.Claims,
		Identifier: request.CompanyIdentifier,
	})
	if err != nil {
		err = exception.Delete{Reasons: []string{"retrieve company error", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// collect any users in the company party
	companyUserCollectResponse, err := a.userRecordHandler.Collect(&userRecordHandler.CollectRequest{
		Claims: a.systemClaims, // using system claims since only system can see users from another party
		Criteria: []criterion.Criterion{
			exactTextCriterion.Criterion{
				Field: "partyId.id",
				Text:  companyRetrieveResponse.Company.Id,
			},
		},
	})
	if err != nil {
		err = exception.Delete{Reasons: []string{"collect users error", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	// delete all users in the company party
	for idx := range companyUserCollectResponse.Records {
		if _, err := a.userRecordHandler.Delete(&userRecordHandler.DeleteRequest{
			Claims: a.systemClaims, // using system claims since only system can see users from another party
			Identifier: id.Identifier{
				Id: companyUserCollectResponse.Records[idx].Id,
			},
		}); err != nil {
			err = exception.Delete{Reasons: []string{"delete company user error", err.Error()}}
			log.Error(err.Error())
			return nil, err
		}
	}

	// delete company
	if _, err := a.companyRecordHandler.Delete(&recordHandler.DeleteRequest{
		Claims:     request.Claims,
		Identifier: request.CompanyIdentifier,
	}); err != nil {
		err = exception.Delete{Reasons: []string{"delete error", err.Error()}}
		log.Error(err.Error())
		return nil, err
	}

	return &administrator2.DeleteResponse{}, nil
}
