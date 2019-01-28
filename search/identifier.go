package search

import "gitlab.com/iotTracker/brain/search/identifier"

type Identifier interface {
	IsValid() error                // Returns the validity of the Identifier
	Type() identifier.Type                    // Returns the IdentifierType of the Identifier
	ToMap() map[string]interface{} // Returns a map to use to query the database
	//Marshal() MarshalledIdentifier // Returns the Identifier in Marshalled form
}
