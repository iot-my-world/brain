package criterion

type Criterion interface {
	IsValid() error                   // Returns the validity of the Criterion
	Type() Type                       // Returns the Type of the Criterion
	ToFilter() map[string]interface{} // Returns a map filter to use to query the databases
}
