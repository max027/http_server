package server

var status_code map[int]string = map[int]string{
	200: "OK",
	204: "No Content",
	404: "Not Found",
	501: "Not Implemented",
	500: "Internal Server Error",
	400: "Bad Request",
	201: "Created",
}
var headers map[string]string = map[string]string{
	"Server":         "CrudeServer",
	"Content-Type":   "text/html",
	"Content-Length": "0",
}
