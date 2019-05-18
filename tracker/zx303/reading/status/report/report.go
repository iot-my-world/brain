package report

type Battery struct {
	Readings [][]int64 `json:"readings"`
}

func NewReadingEntry(timestamp, percentage int64) []int64 {
	return []int64{
		timestamp, percentage,
	}
}
