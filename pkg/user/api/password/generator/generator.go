package generator

type Generator interface {
	Generate(request *GenerateRequest) (*GenerateResponse, error)
}

type GenerateRequest struct {
	CryptoBytesLength int
}

type GenerateResponse struct {
	Password string
}
