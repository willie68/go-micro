package grpcapi

import (
	"context"
	"errors"

	"github.com/willie68/go-micro/pkg/protobuf"
	"google.golang.org/grpc"
)

type APIgRPC struct {
	protobuf.UnimplementedConfigServer
}

var _ protobuf.ConfigServer = &APIgRPC{}

func NewAPIgRPC() (*APIgRPC, error) {
	return &APIgRPC{}, nil
}

func (a *APIgRPC) RegisterGRPC(srv *grpc.Server) {
	protobuf.RegisterConfigServer(srv, a)
}

func (a *APIgRPC) List(context.Context, *protobuf.ListRequest) (*protobuf.ListReply, error) {
	return nil, errors.ErrUnsupported
}
