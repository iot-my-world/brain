package exception

type ReadingCreation struct {
	Reason string
}

func (e ReadingCreation) Error() string {
	return "error creating reading: " + e.Reason
}

type BulkReadingCreation struct {
	Reason string
}

func (e BulkReadingCreation) Error() string {
	return "error creating bulk of readings: " + e.Reason
}
