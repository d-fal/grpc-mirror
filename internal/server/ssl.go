package server

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (srv *Server) setupSSL() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	if srv.tlsEnable { // we dont need tls

		// certifications should be
		creds, sslErr := credentials.NewServerTLSFromFile("cert.crt", "cert.key")

		if sslErr != nil {
			log.Fatal("Cannot open certs ", sslErr)
		}
		opts = append(opts, grpc.Creds(creds))

	}

	return opts
}
