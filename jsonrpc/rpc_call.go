package jsonrpc

import (
	"bytes"
	"net/http"

	"github.com/gorilla/rpc/json"
)

const URL = "http://localhost:5000/rpc"

func RPCCall(method string, params interface{}, out interface{}) error {
	msg, err := json.EncodeClientRequest(method, params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.DecodeClientResponse(resp.Body, out)
	if err != nil {
		return err
	}

	return nil
}
