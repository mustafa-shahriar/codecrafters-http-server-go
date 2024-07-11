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
		if err != nil {
			continue
		}
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)
	if err != nil {
		writeToConn(conn, "400 bad", "", "")
		return
	}

	request := newResquest(buffer)

	switch request.target {
	case "/":
		writeToConn(conn, "200 OK", "", "")
	default:
		writeToConn(conn, "404 Not Found", "", "")
	}

}

func writeToConn(conn net.Conn, statusCode, header, body string) {
	res := fmt.Sprintf("HTTP/1.1 %s\r\n%s\r\n%s", statusCode, header, body)
	conn.Write([]byte(res))
}
