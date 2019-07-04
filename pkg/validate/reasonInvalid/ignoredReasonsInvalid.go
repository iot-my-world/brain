package reasonInvalid

type IgnoredReasonsInvalid struct {
	ReasonsInvalid map[string][]Type
}

func (i IgnoredReasonsInvalid) CanIgnore(reason ReasonInvalid) bool {
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
