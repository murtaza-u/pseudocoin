package jsonrpc

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/murtaza-udaipurwala/pseudocoin/core"
)

var blockchain core.Blockchain
var utxoset core.UTXOSet

type RPC struct{}

func InitRPCServer(dbFile string, bc core.Blockchain, UTXOSet core.UTXOSet) error {
	blockchain = bc
	utxoset = UTXOSet

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(RPC), "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}

	log.Printf("Listening on port :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}
