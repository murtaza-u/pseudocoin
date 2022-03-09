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

type RPC struct{}

var dbFile string

func getBlockchain() (*core.Blockchain, error) {
	bc, err := core.NewBlockchain(dbFile)
	return &bc, err
}

func InitRPCServer(file string) error {
	dbFile = file

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
