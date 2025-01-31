package server

import (
	"bytes"
	"strings"
)

type Parser struct {
	method       string
	uri          string
	http_version string
}

func (par *Parser) parse(data string) {
	lines := strings.Split(data, "\r\n")
	if len(lines) > 0 {
		request_line := []byte(lines[0])

		words := bytes.Split(request_line, []byte(" "))
		par.method = string(words[0])

		if len(words) > 1 {
			par.uri = string(words[1])
		}

		if len(words) > 2 {
			par.http_version = string(words[2])
		}
	}
}
