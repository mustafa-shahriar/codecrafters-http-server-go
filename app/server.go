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
		writeToConn(conn, nil, "400 bad", "", nil)
		return
	}

	request := newResquest(buffer)

	if request.target == "/" {
		writeToConn(conn, request, "200 OK", "", nil)
		return

	}

	if strings.HasPrefix(request.target, "/echo/") {
		path := strings.Split(request.target, "/")
		pathParam := path[len(path)-1]
		body := []byte(pathParam)
		header := fmt.Sprintf("Content-Type: text/plain\r\n")
		writeToConn(conn, request, "200 OK", header, body)
		return
	}

	if request.target == "/user-agent" {
		userAgent := request.header["User-Agent"]
		body := []byte(userAgent)
		header := fmt.Sprintf("Content-Type: text/plain\r\n")
		writeToConn(conn, request, "200 OK", header, body)
		return

	}

	if strings.HasPrefix(request.target, "/files/") && request.method == "GET" {
		path := strings.Split(request.target, "/")
		pathParam := strings.TrimSpace(path[len(path)-1])

		file, err := os.Open(os.Args[2] + pathParam)
		if err != nil {
			fmt.Println(err)
			header := "Content-Type: application/octet-stream\r\n"
			writeToConn(conn, request, "404 Not Found", header, nil)
			return
		}
		fileContent := make([]byte, 1024)
		n, _ := file.Read(fileContent)
		fileContent = fileContent[:n]

		header := fmt.Sprintf("Content-Type: application/octet-stream\r\n")
		writeToConn(conn, request, "200 OK", header, fileContent)
		return
	}

	if strings.HasPrefix(request.target, "/files/") && request.method == "POST" {
		path := strings.Split(request.target, "/")
		fileName := path[len(path)-1]

		file, err := os.Create(os.Args[2] + fileName)
		if err != nil {
			writeToConn(conn, request, "400 Bad Request", "", nil)
			return
		}

		_, err = file.Write(request.body)
		if err != nil {
			writeToConn(conn, request, "400 Bad Request", "", nil)
			return
		}

		writeToConn(conn, request, "201 Created", "", nil)
		return
	}

	writeToConn(conn, request, "404 Not Found", "", nil)
}

func writeToConn(conn net.Conn, request *Request, statusCode, header string, body []byte) {
	if doesAcceptGzip(request) {
		header += "Content-Encoding: gzip\r\n"
		body = compress(body)
	}
	header += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	res := fmt.Sprintf("HTTP/1.1 %s\r\n%s\r\n", statusCode, header)
	resByte := []byte(res)
	resByte = append(resByte, body...)
	conn.Write(resByte)
}
