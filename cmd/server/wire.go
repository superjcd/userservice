//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/superjcd/userservice/config"
	v1 "github.com/superjcd/userservice/genproto/v1"
	"github.com/superjcd/userservice/service"
)

// InitServer Inject service's component
func InitServer(conf *config.Config) (v1.UserServiceServer, error) {

	wire.Build(
		service.NewClient,
		service.NewServer,
	)

	return &service.Server{}, nil

}
