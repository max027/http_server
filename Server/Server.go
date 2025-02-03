package server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

type Server struct {
	Host string
	Port string
}

var blank_line []byte = []byte("\r\n")

func (ser *Server) Start() {
	ser.Host = "tcp"
	ser.Port = ":8888"
	// move to seprate file
	listner, err := net.Listen(ser.Host, ser.Port)
	if err != nil {
		fmt.Println("Error occure while starting server:", err)
		return
	}
	defer listner.Close()
	fmt.Println("Server is listening on port ", ser.Port)
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		ser.Handel_request(conn)
	}
}

func (ser *Server) response_line(status int) []byte {
	reason := status_code[status]

	line := []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, reason))
	return line
}

func (ser *Server) response_header() []byte {
	// extra header
	header := ""
	for key, value := range headers {
		header += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	return []byte(header)
}

func (ser *Server) Handel_request(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	parser := &parser{}
	parser.parse(reader)

	var response []byte
	if parser.method == "GET" {
		response = ser.handel_GET(parser)
	} else if parser.method == "POST" {
		response = ser.handel_POST()
	} else if parser.method == "DELETE" {
		response = ser.handel_DELETE(parser)
	} else if parser.method == "PUT" {
		response = ser.handel_PUT(parser)
	} else {
		response = ser.handel_501()
	}
	conn.Write(response)

}

func (ser *Server) handel_GET(par *parser) []byte {
	filename := strings.Trim(par.uri, "/")
	_, err := os.Stat(fmt.Sprintf("C:\\Users\\saurabh\\programming\\golang\\http_server\\Server\\%s", filename))
	if err != nil {
		response_line := ser.response_line(404)
		response_header := ser.response_header()
		response_body := []byte("<h1>File not found</h1>")
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}

	response_line := ser.response_line(200)
	response_header := ser.response_header()
	response_body, err := readFileAsBytes(fmt.Sprintf("C:\\Users\\saurabh\\programming\\golang\\http_server\\Server\\%s", filename))
	if err != nil {
		fmt.Println("error reading file")
	}

	return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
}

func readFileAsBytes(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ser *Server) handel_POST() []byte {
	var response []byte
	response_line := ser.response_line(201)
	response_header := ser.response_header()
	response_body := []byte("resource created")
	response = bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)

	return response
}

func (ser *Server) handel_501() []byte {
	response_line := ser.response_line(501)

	response_header := ser.response_header()

	response_body := []byte("<h1>501 not Implemented</h1>")

	return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
}

func (ser *Server) handel_DELETE(par *parser) []byte {
	filename := par.uri
	_, err := os.Stat(fmt.Sprintf("C:\\Users\\saurabh\\programming\\golang\\http_server\\Server\\%s", filename))
	response_header := ser.response_header()
	if err != nil {
		fmt.Println("Resource Not found")
		response_body := []byte("<h1>File not found</h1>")
		response_line := ser.response_line(404)
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	response_line := ser.response_line(204)
	err2 := os.Remove(fmt.Sprintf("C:\\Users\\saurabh\\programming\\golang\\http_server\\Server\\%s", filename))
	if err2 != nil {
		fmt.Println("Failed to delete resource")
		response_line := ser.response_line(500)
		response_body := []byte("<h1>Internal Server Error</h1>")
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	return bytes.Join([][]byte{response_line, response_header, blank_line}, nil)
}
func (ser *Server) handel_PUT(par *parser) []byte {
	filename := par.uri
	_, err := os.Stat(fmt.Sprintf("C:\\Users\\saurabh\\programming\\golang\\http_server\\Server\\%s", filename))
	response_header := ser.response_header()
	if err != nil {
		fmt.Println("Resource Not found")
		response_body := []byte("<h1>File not found</h1>")
		response_line := ser.response_line(404)
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	response_line := ser.response_line(201)
	return bytes.Join([][]byte{response_line, response_header, blank_line}, nil)
}
