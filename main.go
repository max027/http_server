package main

import server "github.com/max027/http_server/Server"

func main() {
	server := &server.Server{}
	server.Host = "tcp"
	server.Port = ":8888"
	server.Start()
}
