package genRecordHandler

// Party is a type which returns party details
type Party interface {
	Details() Details
	SetId()
}

