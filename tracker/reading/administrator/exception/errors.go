package exception

type ReadingCreation struct {
	Reason string
}

func (e ReadingCreation) Error() string {
	return "error creating reading: " + e.Reason
}
