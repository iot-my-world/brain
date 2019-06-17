package jsonRpc

import (
	"github.com/iot-my-world/brain/barcode"
	barcodeScanner "github.com/iot-my-world/brain/barcode/scanner"
	"net/http"
)

type adaptor struct {
	barcodeScanner *barcodeScanner.Scanner
}

func New(
	barcodeScanner *barcodeScanner.Scanner,
) *adaptor {
	return &adaptor{
		barcodeScanner: barcodeScanner,
	}
}

type ScanRequest struct {
	ImageData string `json:"imageData"`
}

type ScanResponse struct {
	Barcode barcode.Barcode `json:"barcode"`
}

func (a *adaptor) Scan(r *http.Request, request *ScanRequest, response *ScanResponse) error {
	scanResponse, err := a.barcodeScanner.Scan(&barcodeScanner.ScanRequest{
		ImageData: request.ImageData,
	})
	if err != nil {
		return err
	}

	response.Barcode = scanResponse.Barcode
	return nil
}
