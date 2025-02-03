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
	Host        string
	Port        string
	status_code map[int]string
	headers     map[string]string
}

func (ser *Server) Start() {
	// move to seprate file
	ser.status_code = map[int]string{
		200: "OK",
		204: "No Content",
		404: "Not Found",
		501: "Not Implemented",
		500: "Internal Server Error",
		400: "Bad Request",
		201: "Created",
	}
	ser.headers = map[string]string{
		"Server":         "CrudeServer",
		"Content-Type":   "text/html",
		"Content-Length": "0",
	}

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

func (ser *Server) response_line(status_code int) []byte {
	reason := ser.status_code[status_code]

	line := []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status_code, reason))
	return line
}

func (ser *Server) response_header() []byte {
	// extra header
	header := ""
	for key, value := range ser.headers {
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
		blank_line := []byte("\r\n")
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}

	response_line := ser.response_line(200)
	response_header := ser.response_header()
	blank_line := []byte("\r\n")
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
	blank_line := []byte("\r\n")
	response = bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)

	return response
}

func (ser *Server) handel_501() []byte {
	response_line := ser.response_line(501)

	response_header := ser.response_header()
	blank_line := []byte("\r\n")

	response_body := []byte("<h1>501 not Implemented</h1>")

	return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
}

func (ser *Server) handel_DELETE(par *parser) []byte {
	filename := par.uri
	_, err := os.Stat(fmt.Sprintf("C:\\Users\\saurabh\\programming\\golang\\http_server\\Server\\%s", filename))
	response_header := ser.response_header()
	blank_line := []byte("\r\n")
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
	blank_line := []byte("\r\n")
	if err != nil {
		fmt.Println("Resource Not found")
		response_body := []byte("<h1>File not found</h1>")
		response_line := ser.response_line(404)
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	response_line := ser.response_line(201)
	return bytes.Join([][]byte{response_line, response_header, blank_line}, nil)
}
