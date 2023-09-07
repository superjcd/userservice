//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/HooYa-Bigdata/microservices/grpc_service/userservice/config"
	v1 "github.com/HooYa-Bigdata/microservices/grpc_service/userservice/genproto/v1"
	"github.com/HooYa-Bigdata/microservices/grpc_service/userservice/service"
)

// InitServer Inject service's component
func InitServer(conf *config.Config) (v1.UserServiceServer, error) {

	wire.Build(
		service.NewClient,
		service.NewServer,
	)

	return &service.Server{}, nil

}
