package wrapped

import (
	"encoding/json"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/security/claims"
	apiUserLogin "gitlab.com/iotTracker/brain/security/claims/login/user/api"
	humanUserLogin "gitlab.com/iotTracker/brain/security/claims/login/user/human"
	"gitlab.com/iotTracker/brain/security/claims/registerClientAdminUser"
	"gitlab.com/iotTracker/brain/security/claims/registerClientUser"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyUser"
	"gitlab.com/iotTracker/brain/security/claims/resetPassword"
	wrappedClaimsException "gitlab.com/iotTracker/brain/security/claims/wrapped/exception"
	"net/http"
)

type Wrapped struct {
	Type  claims.Type     `json:"type"`
	Value json.RawMessage `json:"value"`
}

func Wrap(claimsToWrap claims.Claims) (Wrapped, error) {
	if claimsToWrap == nil {
		return Wrapped{}, wrappedClaimsException.Invalid{Reasons: []string{"nil claimsToWrap provided"}}
	}

	marshalledValue, err := json.Marshal(claimsToWrap)
	if err != nil {
		return Wrapped{}, wrappedClaimsException.Wrapping{Reasons: []string{"marshalling", err.Error()}}
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
		var unmarshalledClaims humanUserLogin.Login
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, wrappedClaimsException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.APIUserLogin:
		var unmarshalledClaims apiUserLogin.Login
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, wrappedClaimsException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterCompanyAdminUser:
		var unmarshalledClaims registerCompanyAdminUser.RegisterCompanyAdminUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, wrappedClaimsException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterCompanyUser:
		var unmarshalledClaims registerCompanyUser.RegisterCompanyUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, wrappedClaimsException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterClientAdminUser:
		var unmarshalledClaims registerClientAdminUser.RegisterClientAdminUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, wrappedClaimsException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.RegisterClientUser:
		var unmarshalledClaims registerClientUser.RegisterClientUser
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, wrappedClaimsException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	case claims.ResetPassword:
		var unmarshalledClaims resetPassword.ResetPassword
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, wrappedClaimsException.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	default:
		return nil, wrappedClaimsException.Invalid{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, brainException.Unexpected{Reasons: []string{"identifier still nil"}}
	}

	// check for expiry
	if result.Expired() {
		return nil, wrappedClaimsException.Invalid{Reasons: []string{"expired"}}
	}

	return result, nil
}

func UnwrapClaimsFromContext(r *http.Request) (claims.Claims, error) {
	wrapped, ok := r.Context().Value("wrappedClaims").(Wrapped)
	if !ok {
		return nil, wrappedClaimsException.CouldNotParseFromContext{}
	}
	return wrapped.Unwrap()
}
