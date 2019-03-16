package basic

import (
	brainException "gitlab.com/iotTracker/brain/exception"
	clientAction "gitlab.com/iotTracker/brain/party/client/action"
	clientRecordHandler "gitlab.com/iotTracker/brain/party/client/recordHandler"
	clientRecordHandlerException "gitlab.com/iotTracker/brain/party/client/recordHandler/exception"
	clientValidator "gitlab.com/iotTracker/brain/party/client/validator"
	userRecordHandler "gitlab.com/iotTracker/brain/party/user/recordHandler"
	userRecordHandlerException "gitlab.com/iotTracker/brain/party/user/recordHandler/exception"
	"gitlab.com/iotTracker/brain/search/identifier/adminEmailAddress"
	"gitlab.com/iotTracker/brain/search/identifier/emailAddress"
	"gitlab.com/iotTracker/brain/validate/reasonInvalid"
)

type validator struct {
	clientRecordHandler clientRecordHandler.RecordHandler
	userRecordHandler   userRecordHandler.RecordHandler
}

func New(
	clientRecordHandler clientRecordHandler.RecordHandler,
	userRecordHandler userRecordHandler.RecordHandler,
) clientValidator.Validator {
	return &validator{
		clientRecordHandler: clientRecordHandler,
		userRecordHandler:   userRecordHandler,
	}
}

func (v *validator) ValidateValidateRequest(request *clientValidator.ValidateRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.Claims == nil {
		reasonsInvalid = append(reasonsInvalid, "claims are nil")
	}

	if len(reasonsInvalid) > 0 {
		return brainException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (v *validator) Validate(request *clientValidator.ValidateRequest, response *clientValidator.ValidateResponse) error {
	if err := v.ValidateValidateRequest(request); err != nil {
		return err
	}

	allReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)
	clientToValidate := &request.Client

	if (*clientToValidate).Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "id",
			Type:  reasonInvalid.Blank,
			Help:  "id cannot be blank",
			Data:  (*clientToValidate).Id,
		})
	}

	if (*clientToValidate).Name == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "name",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).Name,
		})
	}

	if (*clientToValidate).ParentPartyType == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentPartyType",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).ParentPartyType,
		})
	}

	if (*clientToValidate).ParentId.Id == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "parentId",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).ParentId,
		})
	}

	if (*clientToValidate).AdminEmailAddress == "" {
		allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
			Field: "adminEmailAddress",
			Type:  reasonInvalid.Blank,
			Help:  "cannot be blank",
			Data:  (*clientToValidate).AdminEmailAddress,
		})
	}

	returnedReasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	// Perform additional checks/ignores considering method field
	switch request.Action {
	case clientAction.Create:

		if (*clientToValidate).AdminEmailAddress != "" {

			// Check if there is another client that is already using the same admin email address
			if err := v.clientRecordHandler.Retrieve(&clientRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: adminEmailAddress.Identifier{
					AdminEmailAddress: (*clientToValidate).AdminEmailAddress,
				},
			},
				&clientRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case clientRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*clientToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*clientToValidate).AdminEmailAddress,
				})
			}

			// Check if there is another user that is already using the same admin email address
			if err := v.userRecordHandler.Retrieve(&userRecordHandler.RetrieveRequest{
				Claims: request.Claims,
				Identifier: emailAddress.Identifier{
					EmailAddress: (*clientToValidate).AdminEmailAddress,
				},
			},
				&userRecordHandler.RetrieveResponse{}); err != nil {
				switch err.(type) {
				case userRecordHandlerException.NotFound:
					// this is what we want, do nothing
				default:
					allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
						Field: "adminEmailAddress",
						Type:  reasonInvalid.Unknown,
						Help:  "unknown error",
						Data:  (*clientToValidate).AdminEmailAddress,
					})
				}
			} else {
				// there was no error, this email is already in database
				allReasonsInvalid = append(allReasonsInvalid, reasonInvalid.ReasonInvalid{
					Field: "adminEmailAddress",
					Type:  reasonInvalid.Duplicate,
					Help:  "already exists",
					Data:  (*clientToValidate).AdminEmailAddress,
				})
			}
		}

		// Check if there is another user that is already using the same admin email address

		// Ignore reasons not applicable for this method
		for _, reason := range allReasonsInvalid {
			if !v.createIgnoredReasons.CanIgnore(reason) {
				returnedReasonsInvalid = append(returnedReasonsInvalid, reason)
			}
		}
	default:
		returnedReasonsInvalid = allReasonsInvalid
	}

	response.ReasonsInvalid = returnedReasonsInvalid
	return nil
}
