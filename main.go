package main

import (
	"bufio"
	"fmt"
	"net"
)

func handelConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error Reading:", err)
		return
	}
	fmt.Printf("recieved message:%s", message)

	_, err = conn.Write([]byte("Message received.\n"))
	if err != nil {
		fmt.Println("Error writing:", err)
	}
}

func main() {
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listning:", err)
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

		handelConnection(conn)
	}
}
