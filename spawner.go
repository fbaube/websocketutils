package websocketutils

import (
        "net"
	"net/http"
	"fmt"
	WS "github.com/coder/websocket"
	// "nhooyr.io/websocket/wsjson"
	// "crypto/tls"
)

type UpgradedWebsocketHandler func(*WSUConn) 

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
	wsuCn.NC = WS.NetConn(ZeroTimeCtx, wsCn, WS.MessageText)
	var ncLA, ncRA net.Addr
	ncLA = wsuCn.NC.LocalAddr()
	ncRA = wsuCn.NC.RemoteAddr()
	wsuCn.LocalAddr  = ncLA
	wsuCn.RemoteAddr = ncRA
	fmt.Printf("Spawned: local<%s:%s> remote<%s:%s> \n",
		ncLA.Network(), ncLA.String(), ncRA.Network(), ncRA.String())

	go GlobalUpgradedWebsocketHandler(wsuCn) 
}

