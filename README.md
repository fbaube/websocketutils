# websocketutils

This compiles Golang to wasm (webassembly). It opens a websocket
back to the server, and then says so on the browser page.

Caveat: It has been tested only on localhost. 

```
GOARCH=wasm GOOS=js go build -o client.wasm client.go
go run ../server/main.go 
```

Then open http://localhost:9090 

The websocket is requested via a call to http://localhost:9090/ws
(If you try to load this in a browser, you get weird errors!)

This was developed initially in Google Chrome on macOS using go1.16,
and  more recently in Safari on macOS using go1.21.4 

