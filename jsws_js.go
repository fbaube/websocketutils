// +build js,wasm

package websocketutils

import (
	"fmt"
	"errors"
	"syscall/js"
)

func TryJSWS(address string) (interface{}, error) {
     	println("trying TryJSWS()...")
	ws := js.Global().Get("WebSocket").New(address)
	if ws.Truthy() {
		fmt.Printf("JS WS OK :: (%T) %+v \n", ws, ws)
		return ws, nil
	}
	return nil, errors.New("TryJSWS FAILED")
}

