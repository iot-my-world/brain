package wrapped

import (
	"encoding/json"
	brainException "github.com/iot-my-world/brain/internal/exception"
	claims2 "github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/api"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	registerClientAdminUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	registerClientUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerClientUser"
	registerCompanyAdminUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	registerCompanyUser2 "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyUser"
	resetPassword2 "github.com/iot-my-world/brain/pkg/security/claims/resetPassword"
	"github.com/iot-my-world/brain/pkg/security/claims/wrapped/exception"
	"net/http"
)

type Wrapped struct {
	Type  claims2.Type    `json:"type"`
	Value json.RawMessage `json:"value"`
}

func Wrap(claimsToWrap claims2.Claims) (Wrapped, error) {
	if claimsToWrap == nil {
		return Wrapped{}, exception.Invalid{Reasons: []string{"nil claimsToWrap provided"}}
	}

	marshalledValue, err := json.Marshal(claimsToWrap)
	if err != nil {
		return Wrapped{}, exception.Wrapping{Reasons: []string{"marshalling", err.Error()}}
	}
	return Wrapped{
		Type:  claimsToWrap.Type(),
		Value: marshalledValue,
	}, nil
}

func (wc Wrapped) Unwrap() (claims2.Claims, error) {
	var result claims2.Claims = nil

	switch wc.Type {
	case claims2.HumanUserLogin:
		var unmarshalledClaims human.Login
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims2.APIUserLogin:
		var unmarshalledClaims api.Login
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims2.RegisterCompanyAdminUser:
		var unmarshalledClaims registerCompanyAdminUser2.RegisterCompanyAdminUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims2.RegisterCompanyUser:
		var unmarshalledClaims registerCompanyUser2.RegisterCompanyUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims2.RegisterClientAdminUser:
		var unmarshalledClaims registerClientAdminUser2.RegisterClientAdminUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims2.RegisterClientUser:
		var unmarshalledClaims registerClientUser2.RegisterClientUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims2.ResetPassword:
		var unmarshalledClaims resetPassword2.ResetPassword
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	default:
		return nil, exception.Invalid{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, brainException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	// check for expiry
	if result.Expired() {
		return nil, exception.Invalid{Reasons: []string{"expired"}}
	}

	return result, nil
}

func UnwrapClaimsFromContext(r *http.Request) (claims2.Claims, error) {
	wrapped, ok := r.Context().Value("wrappedClaims").(Wrapped)
	if !ok {
		return nil, exception.CouldNotParseFromContext{}
	}
	return wrapped.Unwrap()
}
