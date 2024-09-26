// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: order.proto

package order

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "go-micro.dev/v5/client"
	server "go-micro.dev/v5/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Order service

type OrderService interface {
	CheckAndTransformData(ctx context.Context, in *CheckAndTransformDataRequest, opts ...client.CallOption) (*CheckAndTransformDataResponse, error)
}

type orderService struct {
	c    client.Client
	name string
}

func NewOrderService(name string, c client.Client) OrderService {
	return &orderService{
		c:    c,
		name: name,
	}
}

func (c *orderService) CheckAndTransformData(ctx context.Context, in *CheckAndTransformDataRequest, opts ...client.CallOption) (*CheckAndTransformDataResponse, error) {
	req := c.c.NewRequest(c.name, "Order.CheckAndTransformData", in)
	out := new(CheckAndTransformDataResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Order service

type OrderHandler interface {
	CheckAndTransformData(context.Context, *CheckAndTransformDataRequest, *CheckAndTransformDataResponse) error
}

func RegisterOrderHandler(s server.Server, hdlr OrderHandler, opts ...server.HandlerOption) error {
	type order interface {
		CheckAndTransformData(ctx context.Context, in *CheckAndTransformDataRequest, out *CheckAndTransformDataResponse) error
	}
	type Order struct {
		order
	}
	h := &orderHandler{hdlr}
	return s.Handle(s.NewHandler(&Order{h}, opts...))
}

type orderHandler struct {
	OrderHandler
}

func (h *orderHandler) CheckAndTransformData(ctx context.Context, in *CheckAndTransformDataRequest, out *CheckAndTransformDataResponse) error {
	return h.OrderHandler.CheckAndTransformData(ctx, in, out)
}
