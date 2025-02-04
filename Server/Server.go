package server

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Server struct {
	Host string
	Port string
}

var blank_line []byte = []byte("\r\n")

func (ser *Server) Start() {
	ser.Host = "tcp"
	ser.Port = ":8888"
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env ", "ERR", err)
		return
	}

	listner, err := net.Listen(ser.Host, ser.Port)
	if err != nil {
		slog.Error("Error occure while starting server", "Err:", err)
		return
	}

	defer listner.Close()

	slog.Info("Server is Listening", "port", "8888")
	for {
		conn, err := listner.Accept()
		if err != nil {
			slog.Error("Error accepting connection", "ERR", err)
			continue
		}
		go ser.Handel_request(conn)
	}
}

func (ser *Server) response_line(status int) []byte {
	reason := status_code[status]

	line := []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, reason))
	return line
}

func (ser *Server) response_header() []byte {
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
	slog.Info(parser.method, "Path", parser.uri)
	var response []byte
	switch parser.method {
	case "GET":
		response = ser.handel_GET(parser)
	case "POST":
		response = ser.handel_POST(parser)
	case "PUT":
		response = ser.handel_PUT(parser)
	case "DELETE":
		response = ser.handel_DELETE(parser)
	default:
		response = ser.handel_501()
	}
	conn.Write(response)
}

func (ser *Server) handel_GET(par *parser) []byte {
	filename := strings.Trim(par.uri, "/")
	_, err := os.Stat(fmt.Sprintf("%s%s", os.Getenv("FILE_PATH"), filename))
	if err != nil {
		response_line := ser.response_line(404)
		response_header := ser.response_header()
		response_body := []byte("<h1>File not found</h1>")
		slog.Error("File Not found", "Name", filename)
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	response_line := ser.response_line(200)
	response_header := ser.response_header()
	response_body, err := readFileAsBytes(fmt.Sprintf("%s%s", os.Getenv("FILE_PATH"), filename))
	if err != nil {
		slog.Error("error reading file", "Name", filename)
	}
	contentLength := len(response_body)
	headers["Content-Length"] = strconv.Itoa(contentLength)
	return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
}

func readFileAsBytes(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ser *Server) handel_POST(par *parser) []byte {
	var response []byte
	response_line := ser.response_line(201)
	response_header := ser.response_header()
	response_body := []byte("resource created")
	slog.Info("Resource Created", "Name", par.uri)
	response = bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)

	return response
}

func (ser *Server) handel_501() []byte {
	response_line := ser.response_line(501)

	response_header := ser.response_header()

	response_body := []byte("<h1>501 not Implemented</h1>")
	slog.Warn("not Implemented")

	return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
}

func (ser *Server) handel_DELETE(par *parser) []byte {
	filename := par.uri
	_, err := os.Stat(fmt.Sprintf("%s%s", os.Getenv("FILE_PATH"), filename))
	response_header := ser.response_header()
	if err != nil {
		response_body := []byte("<h1>File not found</h1>")
		response_line := ser.response_line(404)
		slog.Error("File Not Found", "Name", filename)
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	response_line := ser.response_line(204)
	err2 := os.Remove(fmt.Sprintf("%s%s", os.Getenv("FILE_PATH"), filename))
	if err2 != nil {
		response_line := ser.response_line(500)
		response_body := []byte("<h1>Internal Server Error</h1>")
		slog.Error("Failed to Delete", "name", filename)
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	slog.Warn("Deleted", "resource", filename)
	return bytes.Join([][]byte{response_line, response_header, blank_line}, nil)
}

func (ser *Server) handel_PUT(par *parser) []byte {
	filename := par.uri
	_, err := os.Stat(fmt.Sprintf("%s%s", os.Getenv("FILE_PATH"), filename))
	response_header := ser.response_header()
	if err != nil {
		response_body := []byte("<h1>File not found</h1>")
		response_line := ser.response_line(404)
		slog.Error("File Not Found", "Name", filename)
		return bytes.Join([][]byte{response_line, response_header, blank_line, response_body}, nil)
	}
	response_line := ser.response_line(201)
	slog.Info("Updated", "Name", filename)
	return bytes.Join([][]byte{response_line, response_header, blank_line}, nil)
}
