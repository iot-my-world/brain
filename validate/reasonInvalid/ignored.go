package reasonInvalid

import "gitlab.com/iotTracker/brain/validate"

type IgnoredReasons struct {
	ReasonsInvalid map[string][]Type
}

func (i IgnoredReasons) CanIgnore(reason validate.ReasonInvalid) bool {
	for fieldString, reasonInvalidTypes := range i.ReasonsInvalid {
		if fieldString == reason.Field {
			for _, reasonType := range reasonInvalidTypes {
				if reasonType == reason.Type {
					return true
				}
			}
		}
	}
	return false
}