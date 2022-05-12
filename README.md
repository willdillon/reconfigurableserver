# reconfigurableserver 
used to allow easily reloading the server without restarting the program. makes it easy to update certs in a running program.

## use:
initialize the new server, then setup the handlers using the server.Mux().
```go
func SetupHandlers() {
	server.Mux().HandleFunc("/", doIndex)
	server.Mux().HandleFunc("/restart", doRestart)
	server.Mux().HandleFunc("/update", doUpdate)
}

func main() {
    server = featureful_server.NewServer(listenaddr, CertificateFile, KeyFile)
    SetupHandlers()
    go server.Start()
    <-server.ShutdownSignal
}
```