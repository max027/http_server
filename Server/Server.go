package server

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type Server struct {
	Host string
	Port string
}

func (ser *Server) Start() {
	listner, err := net.Listen(ser.Host, ser.Port)
	if err != nil {
		fmt.Println("Error occure while starting server:", err)
		return
	}
	defer listner.Close()
	fmt.Println("Server is listening on port 8080...")
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		ser.Handel_request(conn)
	}
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

	response_line := []byte("HTTP/1.1 200 OK\r\n")

	headers := bytes.Join([][]byte{[]byte("Server: Simple Server\r\n"), []byte("Content-Type: text/html\r\n")}, nil)

	blank_line := []byte("\r\n")

	response_body := []byte(`<html>
	<body>
	<h1>Message Recieved</h1>
	</body>
	</html>`)

	_, err = conn.Write(bytes.Join([][]byte{response_line, headers, blank_line, response_body}, nil))
	if err != nil {
		fmt.Println("Error Writing")
	}

}
