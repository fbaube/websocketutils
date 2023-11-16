// NO NO NO NOT YET +build js,wasm

package websocketutils

import (
	"fmt"
	WS "nhooyr.io/websocket"        // github.com/nhooyr/websocket 
	WSJ "nhooyr.io/websocket/wsjson" // github.com/nhooyr/websocket/wsjson
	"time"
	CTX "context"
)

var ctx CTX.Context
var cancel CTX.CancelFunc 

func init() {
	ctx, cancel = CTX.WithTimeout(CTX.Background(), time.Minute)
	// defer cancel()
}

func WriteAndRead(c *WS.Conn, s string) (string, error) {
	err := WSJ.Write(ctx, c, s)
	if err != nil {
		return "", fmt.Errorf("WriteAndRead: " +
		       "wsConn(%#v).Write(%s) failed: %w", c, s, err)
	}
	var v interface{}
	err = WSJ.Read(ctx, c, &v)
	if err != nil {
		return "", fmt.Errorf("WriteAndRead: " +
		       "wsConn(%#v).(wrote:%s).Read() failed: %w", c, s, err)
	}
	return v.(string), nil
}

func WriteAndReturn(c *WS.Conn, s string) error {
	err := WSJ.Write(ctx, c, s)
	if err != nil {
		return fmt.Errorf("WriteAndReturn: " +
		       "wsConn(%#v).Write(%s) failed: %w", c, s, err)
	}
	return nil
}

