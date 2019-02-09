package reasonInvalid

type ReasonInvalid struct {
	Field string      `json:"field"`
	Type  Type        `json:"type"`
	Help  string      `json:"help"`
	Data  interface{} `json:"data"`
}
