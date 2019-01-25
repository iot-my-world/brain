package cors

import (
	"github.com/gorilla/rpc"
	"net/http"
	"strings"
)

// construct a new CORS codec
func CodecWithCors(corsDomains []string, originalCodec rpc.Codec) rpc.Codec {
	return CorsCodec{corsDomains, originalCodec}
}

type CorsCodecRequest struct {
	CorsDomains            []string
	UnderlyingCodecRequest rpc.CodecRequest
}

//override exactly one method of the underlying anonymous field and delegate to it.
func (ccr CorsCodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}, methodErr error) error {
	w.Header().Set("Access-Control-Allow-Origin", strings.Join(ccr.CorsDomains, " "))
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin")
	return ccr.UnderlyingCodecRequest.WriteResponse(w, reply, methodErr)
}

func (ccr CorsCodecRequest) Method() (string, error) {
	return ccr.UnderlyingCodecRequest.Method()
}

func (ccr CorsCodecRequest) ReadRequest(req interface{}) error {
	return ccr.UnderlyingCodecRequest.ReadRequest(req)
}

type CorsCodec struct {
	CorsDomains     []string
	UnderlyingCodec rpc.Codec
}

//override exactly one method of the underlying anonymous field and delegate to it.
func (cc CorsCodec) NewRequest(req *http.Request) rpc.CodecRequest {
	return CorsCodecRequest{cc.CorsDomains, cc.UnderlyingCodec.NewRequest(req)}
}
