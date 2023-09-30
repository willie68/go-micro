package grpc

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	DoGRPC = "sgrpc"
)

type RegisterGRPC interface {
	Register(grpc grpc.Server)
}

// Config configuration of the gRPC service
type Config struct {
	UseSSL   bool   `yaml:"usessl"`
	Port     int    `yaml:"port"`
	CertFile string `yaml:"certfile"`
	KeyFile  string `yaml:"keyfile"`
}

// SgRPC the gRPC service itself
type SgRPC struct {
	cfn        Config
	Started    bool
	lis        net.Listener
	grpcServer *grpc.Server
}

// NewGRPC creates a new gRPC Service, initialize and register it in th di
func NewGRPC(cfn Config) (*SgRPC, error) {
	g := SgRPC{
		cfn: cfn,
	}

	err := g.init()
	if err != nil {
		return nil, err
	}

	do.ProvideNamedValue[SgRPC](nil, DoGRPC, g)
	return &g, nil
}

// init initialize the service the the given configuration
func (g *SgRPC) init() error {
	g.Started = false

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", g.cfn.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	g.lis = lis
	var opts []grpc.ServerOption
	if g.cfn.UseSSL {

		if _, err := os.Stat(g.cfn.CertFile); errors.Is(err, os.ErrNotExist) {
			return err
		}
		if _, err := os.Stat(g.cfn.KeyFile); errors.Is(err, os.ErrNotExist) {
			return err
		}
		creds, err := credentials.NewServerTLSFromFile(g.cfn.CertFile, g.cfn.KeyFile)
		if err != nil {
			return fmt.Errorf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	g.grpcServer = grpc.NewServer(opts...)
}

// GRPCServer getting the underlying GRPCServer
func (g *SgRPC) GRPCServer() *grpc.Server {
	return g.grpcServer
}

// Start starting the GRPC server
func (g *SgRPC) Start() {
	g.grpcServer.Serve(g.lis)
	g.Started = true
}
