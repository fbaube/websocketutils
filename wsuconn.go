package websocketutils

import (
        "fmt"
	"errors"
	S "strings"
	"net"
	"net/http"
	WS "nhooyr.io/websocket"         // github.com/nhooyr/websocket
	WSJ "nhooyr.io/websocket/wsjson" // github.com/nhooyr/websocket/wsjson
)

// WSUConn is a WebSocket Utils connection, containing both 
// an embedded [websocket.Conn] and an embedded [net.Conn]. 
// This makes the API kind of messy - the same verbs appear 
// in different APIs, with different function signatures.
// Also, we rename one of them so that we don't trigger 
// any weird compiler errors. 
// . 
type WSUConn struct {
     *WS.Conn
     NC net.Conn
     LocalAddr, RemoteAddr net.Addr
     Protocol, Address string 
}

// NewWSUConnForDial sets up fields but does not Dial. Since only
// a call to [websocket.Dial] (or [websocket.Accept]) actually 
// creates an actual websocket, this func does not yet either 
// create or initialise the embedded websocket (the [*WS.Conn]).
// . 
func NewWSUConnForDial(address string) (p *WSUConn, e error) {
     	if address == "" { return nil, errors.New(
	   "NewWSUConnForDial: empty address") } 
        p = new(WSUConn)
	p.Address = address 
        i := S.Index(address, "://")
	if i <= 0 {
	   return nil, errors.New("NewWSUConnForDial: no protocol " +
	   	  "(followed by \"://\") in: " + address)
	}
	p.Protocol = address[0:i]	
     	return p, nil
}

// NewWSUConnFromAccept does a [websocket.Accept] to get 
// a [websocket.Conn] and then embeds it in a [WSConn].
// From the API docs:
//  - Accept(w, req, opts *AcceptOptions) (*Conn, error)
//  - Accept accepts a WebSocket handshake from a client
//    and upgrades the the connection to a WebSocket.
//  - Accept will not allow cross origin requests by default.
//  - Accept will write a response to w on all errors.
//  - NOTE: So what the heck does req do here ?? 
// .
/*
func NewWSUConnFromAccept(w http.ResponseWriter, req *http.Request) (*WSUConn, error) {
     var wsc *WS.Conn
     var e error 
     wsc, e = WS.Accept(w, req, nil) // nil is AcceptOptions 
     if e != nil {
     	var p *WSUConn
        p = new(WSUConn)
        p.Conn = wsc
	return p, nil
     }
     return nil, fmt.Errorf("NewWSConnFromAccept failed: %w", e)
}
*/

// Dial calls the underlying [websocket.Dial] to get an actual websocket.
// From the API docs:
//  - Dial(ctx, u string, opts *DialOpts) (*Conn, *http.Response, error)
//  - Dial performs a WebSocket handshake on url.
//  - The *http.Response is the WebSocket handshake response from
//    the server, but you never need to close the rsp.Body yourself.
//  - If an error occurs, the returned response may be non nil,
//    but you can only read the first 1024 bytes of the body.
//  - http:// & https:// URLs work and are interpreted as ws/wss.
//  - NOTE: *http.Response from Dial is &http.Response{}
//    with a 101 status code on success
//  - NOTE: It specifies a text connection, not a binary connection. 
// . 
func (p *WSUConn) Dial() error {
	var pRsp *http.Response
	var err error
	p.Conn, pRsp, err = WS.Dial(DefaultCtx, p.Address, nil)
	if err == nil {
	       fmt.Printf("Dial: new websocket OK: %#v \n", *p.Conn)
	       // NOTE the use of Text, cos we plan to use net/textproto 
	       p.NC = WS.NetConn(ZeroTimeCtx, p.Conn, WS.MessageText)
	       // How to test whether it is valid ? 
	       return nil	
	}
	fmt.Printf("Dial failed: http.Response: %#v \n", *pRsp)
	/*
	println("Dial: websocket creation failed, " +
		"so now tryng to fetch external JS websocket...")
	v, err := TryJSWS()
	fmt.Printf("TryJSWS ret'd: (%T) %#v (%s) \n", v, v, err.Error())
	return nil, err
	*/
	return fmt.Errorf("Dial: websocket creation failure: %w", err) 
}

func (p *WSUConn) WriteAndRead(s string) (string, error) {
	err := WSJ.Write(DefaultCtx, p.Conn, s)
	if err != nil {
		return "", fmt.Errorf("WriteAndRead: " +
		       "wsJsonConn(%#v).Write(%s) failed: %w", p.Conn, s, err)
	}
	println("SENT:", s) 
	var v interface{}
	err = WSJ.Read(DefaultCtx, p.Conn, &v)
	if err != nil {
	       println("WriteAndRead: Read failed:", err.Error())
	       return "", fmt.Errorf("WriteAndRead: wsJsonConn(%#v)." +
		       "(wrote:%s).Read() failed: %w", p.Conn, s, err)
	}
	s = v.(string)
	println("RCVD:", s)
	return s, nil
}

func (p *WSUConn) WriteAndReturn(s string) error {
	err := WSJ.Write(DefaultCtx, p.Conn, s)
	if err != nil {
		return fmt.Errorf("WriteAndReturn: " +
		       "wsJsonConn(%#v).Write(%s) failed: %w", p.Conn, s, err)
	}
	println("SENT(NRN):", s)
	return nil
}

func (p *WSUConn) Read() (string, error) {
        var e error 
        // Read(ctx, c *websocket.Conn, v interface{}) error
        // Read reads a JSON message from c into v. It will
        // reuse buffers in between calls to avoid allocations.
	var v interface{}
        e = WSJ.Read(DefaultCtx, p.Conn, &v)
        if e != nil {
                fmt.Printf("Read: wsjson.Read: %s \n", e.Error())
		return "", fmt.Errorf("Read: %w", e)
        }
        fmt.Printf("RCVD: (%T) %v \n", v, v)
 	return v.(string), nil
}

