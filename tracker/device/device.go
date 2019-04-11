package device

type Device interface {
	Type() Type
	//MarshalJSON() ([]byte, error)
	//UnmarshalJSON(data []byte) error
}
