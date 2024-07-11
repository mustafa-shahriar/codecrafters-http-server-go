package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	defer l.Close()
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Println("Failed to bind to port 4221")
			os.Exit(1)
		}
		conn.Write([]byte("hello wolrd"))
	}

}
