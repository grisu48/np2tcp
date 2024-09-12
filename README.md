# np2tcp

Bidirectional NamedPipe to TCP adapter.

**THIS IS WIP**, everything here is still fresh out of the oven!

Currently the project consists of three different binaries for Windows only:

* `np2tcp` - The actual NamedPipe to TCP adapter
* `echo` - A simple NamedPipe echo server
* `telnet` - Telnet clone, for testing

## np2tcp

Usage: `np2tcp PIPE BINDADDRESS`, where `PIPE` is the path to the named pipe and `BINDADDRESS` is the remote address identifier where the server should listen to.

Usage examples:

```
np2tcp \\.\pipe\10001 127.0.0.1:10001        # Listens on localhost only
np2tcp \\.\pipe\10001 :10001                 # Listen on all interfaces
```

The server allows only one tcp connection at a time, however it allows a client to reconnect, in case the connection is lost.

If the named pipe is closed, the server terminates.

## Building

Building at the moment needs to happen manually:

```
go build -o telnet.exe cmd/telnet/telnet.go
go build -o echo.exe cmd/echo/echo.go
go build -o np2tcp cmd/np2tcp/main.go
```