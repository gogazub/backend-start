//go:build client

package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("net.Dial() error")
	}

	fmt.Println(bufio.NewReader(conn).ReadString('\n'))

}
