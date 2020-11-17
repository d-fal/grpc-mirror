package server

import (
	"context"
	"log"
)

// Listen listens to grpc ip and port
func (srv *Server) Listen(ctx context.Context) error {

	srv.gctx = ctx
	var err error

	if err = srv.grpcServer.Serve(srv.listener); err != nil {
		log.Printf("Failed to server: %v\n", err)
	}

	return err
}
