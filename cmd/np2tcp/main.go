package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/Microsoft/go-winio"
)

var client net.Conn
var npipe net.Conn
var err error

func terminate(rc int) {
	// Gracefully close
	if client != nil {
		client.Close()
	}
	if npipe != nil {
		npipe.Close()
	}
	os.Exit(rc)
}

// Receiver goroutine from the named pipe. Once the pipe is closed (or errors), the program terminates
func handlePipe() {
	buffer := make([]byte, 2048)
	for {
		n, err := npipe.Read(buffer)

		// If the pipe is closed, terminate normally
		if err != nil {
			switch {
			case errors.Is(err, io.ErrClosedPipe),
				errors.Is(err, io.EOF),
				errors.Is(err, syscall.EPIPE):
				fmt.Fprintf(os.Stderr, "pipe closed.\n")
				terminate(0)
			}
			fmt.Fprintf(os.Stderr, "pipe read error: %s\n", err)
			terminate(1)
		}

		if client == nil || n == 0 {
			continue
		}

		if _, err := client.Write(buffer[:n]); err != nil {
			fmt.Fprintf(os.Stderr, "send error: %s\n", err)
			terminate(1)
		}
	}
}

// Routine for handling the tcp client.
func handleClient() error {
	buffer := make([]byte, 2048)
	for client != nil {
		n, err := client.Read(buffer)

		if err != nil {
			return err
		}

		if n == 0 {
			continue
		}

		if _, err = npipe.Write(buffer[:n]); err != nil {
			return err
		}
	}
	return fmt.Errorf("client not ready")
}

func main() {
	pipePath := os.Args[1]
	bindAddress := "127.0.0.1:10001"

	// Open named pipe. If the named pipe closes, the program terminates
	npipe, err = winio.DialPipe(pipePath, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening pipe: %s\n", err)
		os.Exit(1)
	}
	defer npipe.Close()
	go handlePipe()

	// Setup tcp server
	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listening on %s: %s\n", bindAddress, err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Fprintf(os.Stderr, "Listening on %s\n", bindAddress)

	// Allow one client at a time, but if a client disconnects (or is lost) allow reconnections
	for {
		client, err = listener.Accept()
		if err != nil {
			log.Fatalf("error accepting client: %s", err)
			listener.Close()
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "connected: %s\n", client.RemoteAddr().String())
		if err := handleClient(); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				// Terminate program if the pipe is closed
				fmt.Fprintf(os.Stderr, "pipe closed.")
				terminate(0)
			} else if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "client %s disconnected\n", client.RemoteAddr().String())
			} else {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}
		}
	}
}
