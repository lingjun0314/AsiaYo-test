package main

import (
	"orderService/handler/order"
	pb "orderService/proto/order"

	"github.com/micro/plugins/v5/registry/consul"
	"go-micro.dev/v5"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
)

var (
	service = "order"
	version = "latest"
)

func main() {
	//	Set consul address
	consulReg := consul.NewRegistry(
		registry.Addrs("localhost:8499"),
	)

	//	Create new service
	service := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Registry(consulReg),
	)

	//	Initialise flags
	service.Init()

	//	Register handler
	pb.RegisterOrderHandler(service.Server(), new(order.Order))

	//	Start service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
