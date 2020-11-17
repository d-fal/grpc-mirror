package server

import "grpc-mirror/internal/broker"

// RegisterServer prepares and registers grpc server parameters
func RegisterServer(host string, port uint16) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

// SetTLS sets tls
func (srv *Server) SetTLS(tls bool) *Server {

	srv.tlsEnable = tls

	return srv
}

// SetupBroker setup a broker
func (srv *Server) SetupBroker(broker *broker.Broker) *Server {
	srv.Broker = broker

	return srv
}
