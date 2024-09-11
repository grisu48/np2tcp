package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	remote := os.Args[1]

	conn, err := net.Dial("tcp", remote)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection error: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Fprintf(os.Stderr, "connected: %s\n", remote)

	// Receive
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "receive error: %s\n", err)
			os.Exit(1)
		}
	}()

	// Read from console
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text() + "\n"
		if _, err := conn.Write([]byte(text)); err != nil {
			fmt.Fprintf(os.Stderr, "send error: %s\n", err)
			os.Exit(1)
		}
	}
	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %s\n", err)
		os.Exit(1)
	}
}
