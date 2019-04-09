package scanner

import "gitlab.com/iotTracker/brain/barcode"

type Scanner struct {
}

func New() *Scanner {
	return &Scanner{}
}

type ScanRequest struct {
	ImageData string
}

type ScanResponse struct {
	Barcode barcode.Barcode
}

func (s *Scanner) Scan(request *ScanRequest) (*ScanResponse, error) {
	// perform scanning here with request.ImageData

	return &ScanResponse{
		Barcode: barcode.Barcode{
			Data: "1234",
		},
	}, nil
}
