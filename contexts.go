package websocketutils

import (
	"time"
	CTX "context"
)

var DefaultCtx, ZeroTimeCtx CTX.Context
var DefaultCtxCancel CTX.CancelFunc 

func init() {
	DefaultCtx, DefaultCtxCancel =
		    	 CTX.WithTimeout(CTX.Background(), time.Minute)
	ZeroTimeCtx, _ = CTX.WithTimeout(CTX.Background(), 0)
	// defer cancel()
}

