package main

import server "github.com/max027/http_server/Server"

func main() {
	server := &server.Server{}
	server.Start()
}
