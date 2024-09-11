package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

func main() {
	path := os.Args[1]
	pipe, err := winio.ListenPipe(path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Listening on %s ... \n", path)
	for {
		conn, err := pipe.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept error: %s\n", err)
			os.Exit(1)
		}
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "%s disconnected\n", conn.RemoteAddr().String())
				return
			}
			fmt.Fprintf(os.Stderr, "read error: %s\n", err)
			return
		}
		if _, err := conn.Write(buffer[:n]); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %s\n", err)
			return
		}
		fmt.Printf("echo %d bytes\n", n)
	}

}
