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

var bc core.Blockchain
var UTXOSet core.UTXOSet

type RPC struct{}

func InitRPCServer(dbFile string) error {
	var err error
	bc, err = core.NewBlockchain(dbFile)
	if err != nil {
		return err
	}
	defer bc.DB.Close()

	UTXOSet = core.UTXOSet{Blockchain: &bc}

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(RPC), "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)

	port := os.Getenv("PORT")
	log.Printf("Listening on port :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}
