package identifier

type MarshalledIdentifier string
type IdentifierType string

type Identifier interface {
	IsValid() error // Returns the validity of the Identifier
	Type() IdentifierType // Returns the IdentifierType of the Identifer
	ToMap() map[string]interface{} // Returns a map to use to query the database
	Marshal() MarshalledIdentifier // Returns the Identifier in Marshalled form
}