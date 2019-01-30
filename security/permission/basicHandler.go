package permission

import (
	"gitlab.com/iotTracker/brain/party/user"
	"gitlab.com/iotTracker/brain/security/role"
	"fmt"
	globalException "gitlab.com/iotTracker/brain/exception"
)

type basicHandler struct {
	userRecordHandler user.RecordHandler
	roleRecordHandler role.RecordHandler
}

func NewBasicHandler(
	userRecordHandler user.RecordHandler,
	roleRecordHandler role.RecordHandler,
) *basicHandler {
	return &basicHandler{
		userRecordHandler: userRecordHandler,
		roleRecordHandler: roleRecordHandler,
	}
}

func (bh *basicHandler) ValidateUserHasPermissionRequest(request *UserHasPermissionRequest) error {
	reasonsInvalid := make([]string, 0)

	if request.UserIdentifier == nil {
		reasonsInvalid = append(reasonsInvalid, "identifier is nil")
	} else {
		if !user.IsValidIdentifier(request.UserIdentifier) {
			reasonsInvalid = append(reasonsInvalid, fmt.Sprintf("identifier of type %s not supported for user", request.UserIdentifier.Type()))
		}
	}

	if request.Permission == "" {
		reasonsInvalid = append(reasonsInvalid, "permission is blank")
	}

	if len(reasonsInvalid) > 0 {
		return globalException.RequestInvalid{Reasons: reasonsInvalid}
	} else {
		return nil
	}
}

func (bh *basicHandler) UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error {
	if err := bh.ValidateUserHasPermissionRequest(request); err != nil {
		return err
	}

	return nil
}