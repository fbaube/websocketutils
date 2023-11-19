package websocketutils

import (
	"time"
	CTX "context"
)

var DefaultCtx CTX.Context
var DefaultCtxCancel CTX.CancelFunc 

func init() {
	DefaultCtx, DefaultCtxCancel =
		    CTX.WithTimeout(CTX.Background(), time.Minute)
	// defer cancel()
}

