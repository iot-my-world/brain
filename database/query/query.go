package query

type QueryType string

const ID QueryType = "ID"

type Query struct {
	Type QueryType `json:"queryType"`
	Value string `json:"value"`
}

func (q Query) Parse () {
	switch q.Type {
	case ID:

	}
}

