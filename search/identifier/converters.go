package identifier

type MarshalledIdentifier string

// This function combines the IdentifierType with the JsonMarshalled Identifier itself to form a MarshallIdentifier
// which can then be persisted
func CreateLabelledMarshalledIdentifier(identifierType Type, marshalledIdentifier string) MarshalledIdentifier {
	s := string(identifierType) + "::" + marshalledIdentifier
	return MarshalledIdentifier(s)
}
