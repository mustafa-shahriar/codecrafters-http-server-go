package main

import (
	"bytes"
	"compress/gzip"
	"strings"
)

func compress(buf []byte) []byte {
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	gzipWriter.Write(buf)
	gzipWriter.Close()

	return buffer.Bytes()
}

func doesAcceptGzip(request *Request) bool {
	if request.header["Accept-Encoding"] != "" {
		headers := strings.Split(request.header["Accept-Encoding"], ",")
		for _, h := range headers {
			if strings.TrimSpace(h) == "gzip" {
				return true
			}
		}
	}

	return false
}
