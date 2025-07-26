//go:build server

package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	fmt.Fprintf(conn, "OK\n")
	conn.Close()
}

func startServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("net.Listen() error")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept() error")
		}
		go handleConnection(conn)
	}

}

func main() {
	startServer()
}
