package wrappedClaims

import (
	"encoding/json"
	brainException "gitlab.com/iotTracker/brain/exception"
	"gitlab.com/iotTracker/brain/security/claims"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/security/claims/registerCompanyAdminUser"
	wrappedClaimsException "gitlab.com/iotTracker/brain/security/wrappedClaims/exception"
)

type WrappedClaims struct {
	Type  claims.Type     `json:"type"`
	Value json.RawMessage `json:"value"`
}

func Wrap(claimsToWrap claims.Claims) (WrappedClaims, error) {
	if claimsToWrap == nil {
		return WrappedClaims{}, wrappedClaimsException.Invalid{Reasons: []string{"nil claimsToWrap provided"}}
	}

	marshalledValue, err := json.Marshal(claimsToWrap)
	if err != nil {
		return WrappedClaims{}, wrappedClaimsException.Wrapping{Reasons: []string{"marshalling", err.Error()}}
	}
	return WrappedClaims{
		Type:  claimsToWrap.Type(),
		Value: marshalledValue,
	}, nil
}

func (wc WrappedClaims) Unwrap() (claims.Claims, error) {
	var result claims.Claims = nil

	switch wc.Type {
	case claims.Login:
		var unmarshalledClaims login.Login
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
