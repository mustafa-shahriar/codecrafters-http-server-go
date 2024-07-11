package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var fileDir *string

func main() {
	fileDir = flag.String("directory", "/tmp/", "")
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

	if request.target == "/" {
		writeToConn(conn, "200 OK", "", "")
		return

	}

	if strings.HasPrefix(request.target, "/echo/") {
		path := strings.Split(request.target, "/")
		pathParam := path[len(path)-1]
		header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", len(pathParam))
		writeToConn(conn, "200 OK", header, pathParam)
		return
	}

	if request.target == "/user-agent" {
		userAgent := request.header["User-Agent"]
		header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", len(userAgent))
		writeToConn(conn, "200 OK", header, userAgent)
		return

	}

	if strings.HasPrefix(request.target, "/files/") {
		path := strings.Split(request.target, "/")
		pathParam := path[len(path)-1]

		file, err := os.Open(*fileDir + "/" + pathParam)
		if err != nil {
			header := "Content-Type: application/octet-stream\r\n"
			writeToConn(conn, "404 Not Found", header, "")
			return
		}
		fileContent := make([]byte, 1024)
		n, _ := file.Read(fileContent)
		fileContent = fileContent[:n]

		header := fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d\r\n", n)
		writeToConn(conn, "200 OK", header, string(fileContent))
		return
	}

	writeToConn(conn, "404 Not Found", "", "")
}

func writeToConn(conn net.Conn, statusCode, header, body string) {
	res := fmt.Sprintf("HTTP/1.1 %s\r\n%s\r\n%s", statusCode, header, body)
	conn.Write([]byte(res))
}
