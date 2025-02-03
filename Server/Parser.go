package server

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	method       string
	uri          string
	http_version string
	header       map[string]string
	body         string
}

func (par *parser) parse(reader *bufio.Reader) {
	data, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error Reading:", err)
	}
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
	par.header = make(map[string]string)
	var contentLength int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error Reading header:", err)
		}
		trimmed_line := strings.Trim(line, "\r\n")

		if trimmed_line == "" {
			break
		}
		mp := strings.SplitN(trimmed_line, ":", 2)
		if len(mp) == 2 {
			key := strings.TrimSpace(mp[0])
			value := strings.TrimSpace(mp[1])
			par.header[key] = value
			if key == "Content-Length" {
				contentLength, _ = strconv.Atoi(value)
			}
		}
	}
	if contentLength > 0 {
		temp := make([]byte, contentLength)
		_, err := reader.Read(temp)
		if err != nil {
			fmt.Println("Error Reading body:", err)
			return
		}
		par.body = string(temp)
	} else {
		par.body = ""
	}
	fmt.Println(par.body)
}
