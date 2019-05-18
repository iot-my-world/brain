package generator

type Generator interface {
	Battery(request *BatteryRequest) (*BatteryResponse, error)
}

type BatteryRequest struct {
}

type BatteryResponse struct {
}
