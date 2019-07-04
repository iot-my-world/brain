package environment

type Type string

const Development Type = "development"
const Production Type = "production"

func (t Type) String() string {
	return string(t)
}
