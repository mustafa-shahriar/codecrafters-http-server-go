package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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
		writeToConn(conn, nil, "400 bad", "", "")
		return
	}

	request := newResquest(buffer)

	if request.target == "/" {
		writeToConn(conn, request, "200 OK", "", "")
		return

	}

	if strings.HasPrefix(request.target, "/echo/") {
		path := strings.Split(request.target, "/")
		pathParam := path[len(path)-1]
		header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", len(pathParam))
		writeToConn(conn, request, "200 OK", header, pathParam)
		return
	}

	if request.target == "/user-agent" {
		userAgent := request.header["User-Agent"]
		header := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n", len(userAgent))
		writeToConn(conn, request, "200 OK", header, userAgent)
		return

	}

	if strings.HasPrefix(request.target, "/files/") && request.method == "GET" {
		fmt.Println(request.method)
		path := strings.Split(request.target, "/")
		pathParam := strings.TrimSpace(path[len(path)-1])

		file, err := os.Open(os.Args[2] + pathParam)
		if err != nil {
			fmt.Println(err)
			header := "Content-Type: application/octet-stream\r\n"
			writeToConn(conn, request, "404 Not Found", header, "")
			return
		}
		fileContent := make([]byte, 1024)
		n, _ := file.Read(fileContent)
		fileContent = fileContent[:n]

		header := fmt.Sprintf("Content-Type: application/octet-stream\r\nContent-Length: %d\r\n", n)
		writeToConn(conn, request, "200 OK", header, string(fileContent))
		return
	}

	if strings.HasPrefix(request.target, "/files/") && request.method == "POST" {
		path := strings.Split(request.target, "/")
		fileName := path[len(path)-1]

		file, err := os.Create(os.Args[2] + fileName)
		if err != nil {
			writeToConn(conn, request, "400 Bad Request", "", "")
			return
		}

		_, err = file.Write(request.body)
		if err != nil {
			writeToConn(conn, request, "400 Bad Request", "", "")
			return
		}

		writeToConn(conn, request, "201 Created", "", "")
		return
	}

	writeToConn(conn, request, "404 Not Found", "", "")
}

func writeToConn(conn net.Conn, request *Request, statusCode, header, body string) {
	if request.header["Accept-Encoding"] == "gzip" {
		header += "Content-Encoding: gzip\r\n"
	}
	res := fmt.Sprintf("HTTP/1.1 %s\r\n%s\r\n%s", statusCode, header, body)
	conn.Write([]byte(res))
}
