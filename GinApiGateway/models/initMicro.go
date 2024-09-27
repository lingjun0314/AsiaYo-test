package models

import (
	"github.com/micro/plugins/v5/registry/consul"
	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/registry"
)

var MicroClient client.Client

func init() {
	//	Registry consul
	consulReg := consul.NewRegistry(
		registry.Addrs("172.31.40.181:8500"),
	)

	//	Create service
	service := micro.NewService(
		micro.Registry(consulReg),
	)

	//	Init service
	service.Init()

	//	Set client
	MicroClient = service.Client()
}
