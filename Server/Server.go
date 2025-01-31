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
		404: "Not Found",
		501: "Not Implemented",
	}
	ser.headers = map[string]string{
		"Server":       "CrudeServer",
		"Content-Type": "text/html",
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
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error Reading:", err)
		return
	}
	fmt.Printf("recieved message:%s", message)

	parser := &Parser{}
	parser.parse(message)

	var response_body []byte
	if parser.method == "GET" {
		response_body = ser.handel_GET(parser)
	}

	response_line := ser.response_line(200)

	header := ser.response_header()

	blank_line := []byte("\r\n")

	_, err = conn.Write(bytes.Join([][]byte{response_line, header, blank_line, response_body}, nil))
	if err != nil {
		fmt.Println("Error Writing")
	}
}

func (ser *Server) handel_GET(par *Parser) []byte {
	filename := strings.Trim(par.uri, "/")
	// problem
	path, err := os.Stat(filename)
	fmt.Println("path:", path)
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
	response_body, err := readFileAsBytes(par.uri)
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

func (ser *Server) handel_501() []byte {
	response_line := ser.response_line(501)

	response_header := ser.response_header()
	blank_line := []byte("\r\n")

	response_body := []byte("<h1>501 not Implemented</h1>")

	return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
}
