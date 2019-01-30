package permission

import (
	"gitlab.com/iotTracker/brain/search"
	"gitlab.com/iotTracker/brain/security"
)

type Handler interface {
	UserHasPermission(request *UserHasPermissionRequest, response *UserHasPermissionResponse) error
}

type UserHasPermissionRequest struct {
	UserIdentifier search.Identifier   `json:"userIdentifier"`
	Permission     security.Permission `json:"permission"`
}

type UserHasPermissionResponse struct {
	Result bool `json:"result"`
}
