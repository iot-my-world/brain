package identifier

import "gopkg.in/mgo.v2/bson"

type Identifier interface {
	IsValid() error   // Returns the validity of the Identifier
	Type() Type       // Returns the IdentifierType of the Identifier
	ToFilter() bson.M // Returns a map filter to use to query the databases
}
