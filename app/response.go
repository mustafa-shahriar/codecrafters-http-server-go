package main

import (
	"bytes"
	"strings"
)

type Request struct {
	method string
	target string
	header map[string]string
}

func newResquest(byteArray []byte) *Request {
	r := Request{}

	reader := bytes.NewReader(byteArray)
	var buffer bytes.Buffer
	readline(reader, &buffer)
	rlArray := strings.Split(buffer.String(), " ")
	r.method = rlArray[0]
	r.target = rlArray[1]

	return &r
}

func readline(reader *bytes.Reader, buffer *bytes.Buffer) {
	for {
		b, _ := reader.ReadByte()
		if b == 13 {
			reader.ReadByte()
			break
		}
		buffer.WriteByte(b)
	}

}
