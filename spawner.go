package websocketutils

import (
	"net/http"
	"fmt"
	WS "nhooyr.io/websocket"
	// "nhooyr.io/websocket/wsjson"
	// "crypto/tls"
)

type UpgradedWebsocketHandler func(*WSUConn) 
// func(p *WS.Conn)(*WSU.WSUConn, error)

var GlobalUpgradedWebsocketHandler UpgradedWebsocketHandler

func SpawnUpgradedWSs(w http.ResponseWriter, req *http.Request) { 
	// Accept(w, req, opts *AcceptOptions) (*Conn, error)
	// Accept accepts a WebSocket handshake from a client
	//    and upgrades the the connection to a WebSocket.
	// Accept will not allow cross origin requests by default.
	// Accept will write a response to w on all errors.
	var wsCn *WS.Conn
	var err error 
	wsCn, err = WS.Accept(w, req, nil) // nil is AcceptOptions 
	if err != nil {
		fmt.Printf("Spawner: Accept failed: %s \n", err.Error())
		return 
	}
	
	// "Upgrade"
	wsuCn := new(WSUConn)
	wsuCn.Conn = wsCn
	// Make a net.Conn and get info from it
	wsuCn.NC = WS.NetConn(DefaultCtx, wsCn, WS.MessageText)
	la := wsuCn.NC.LocalAddr()
	ra := wsuCn.NC.RemoteAddr()
	fmt.Printf("Spawned: local<%s:%s> remote<%s:%s> \n",
		la.Network(), la.String(), ra.Network(), ra.String())

	go GlobalUpgradedWebsocketHandler(wsuCn) 
}

