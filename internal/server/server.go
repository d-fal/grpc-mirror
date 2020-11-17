package server

import (
	"context"
	"fmt"
	"grpc-mirror/internal/broker"
	"grpc-mirror/pkg/logger"
	topupService "grpc-mirror/pkg/protobufs/destinationpb"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server struct to initiate grpc server
type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
	gctx       context.Context
	host       string
	port       uint16
	tlsEnable  bool
	Broker     *broker.Broker
	logger     *logger.ApplicationLog
}

// Proxy is the endpoint we are receiving messages on her behalf
type Proxy struct {
	Pub     *grpc.ClientConn
	message chan interface{}
}

// SetupServerConnection to grpc connections
func (srv *Server) SetupServerConnection(f ...interface{}) error {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.host, srv.port))

	if err != nil {
		log.Panicf("Failed to listen %v\n", err)
		return err
	}

	opts := srv.setupSSL()

	s := grpc.NewServer(opts...)

	for _, _f := range f {

		switch _f.(type) {
		case func(*grpc.Server, topupService.TopupServiceServer):
			fn, ok := _f.(func(*grpc.Server, topupService.TopupServiceServer))
			if ok {
				fn(s, srv)
			}
		}

	}
	//enable reflection
	reflection.Register(s)

	srv.listener = listener
	srv.setGRPCServer(s)

	return nil
}

// SetLogger set zap logger
func (srv *Server) SetLogger(l *zap.Logger) {

	_logger := logger.Prepare(l)
	srv.logger = _logger
}

// making the code more readable, we use setters rather than direct assignment
func (srv *Server) setGRPCServer(s *grpc.Server) {
	srv.grpcServer = s
}
