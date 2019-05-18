package report

type Battery struct {
	Readings []BatteryReading `json:"readings"`
}

type BatteryReading struct {
	Percentage int64 `json:"batteryPercentage"`
	Timestamp  int64 `json:"timestamp"`
}
