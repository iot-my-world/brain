package wrapped

import (
	"encoding/json"
	brainException "github.com/iot-my-world/brain/internal/exception"
	"github.com/iot-my-world/brain/pkg/security/claims"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/api"
	"github.com/iot-my-world/brain/pkg/security/claims/login/user/human"
	registerClientAdminUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerClientAdminUser"
	registerClientUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerClientUser"
	registerCompanyAdminUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyAdminUser"
	registerCompanyUserClaims "github.com/iot-my-world/brain/pkg/security/claims/registerCompanyUser"
	resetPasswordClaims "github.com/iot-my-world/brain/pkg/security/claims/resetPassword"
	"github.com/iot-my-world/brain/pkg/security/claims/wrapped/exception"
	"net/http"
)

type Wrapped struct {
	Type  claims.Type     `json:"type"`
	Value json.RawMessage `json:"value"`
}

func Wrap(claimsToWrap claims.Claims) (Wrapped, error) {
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

func (wc Wrapped) Unwrap() (claims.Claims, error) {
	var result claims.Claims = nil

	switch wc.Type {
	case claims.HumanUserLogin:
		var unmarshalledClaims human.Login
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.APIUserLogin:
		var unmarshalledClaims api.Login
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterCompanyAdminUser:
		var unmarshalledClaims registerCompanyAdminUserClaims.RegisterCompanyAdminUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterCompanyUser:
		var unmarshalledClaims registerCompanyUserClaims.RegisterCompanyUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterClientAdminUser:
		var unmarshalledClaims registerClientAdminUserClaims.RegisterClientAdminUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterClientUser:
		var unmarshalledClaims registerClientUserClaims.RegisterClientUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.ResetPassword:
		var unmarshalledClaims resetPasswordClaims.ResetPassword
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

func UnwrapClaimsFromContext(r *http.Request) (claims.Claims, error) {
	wrapped, ok := r.Context().Value("wrappedClaims").(Wrapped)
	if !ok {
		return nil, exception.CouldNotParseFromContext{}
	}
	return wrapped.Unwrap()
}
