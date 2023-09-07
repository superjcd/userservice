package main

import (
	"flag"

	"github.com/HooYa-Bigdata/microservices/grpc_service/userservice/cmd/server"
)

var cfg = flag.String("config", "config/config.yaml", "config file location")

// main main
func main() {
	server.Run(*cfg)
}
