package routing

import (
	"grpc-mirror/pkg/conf"
	"grpc-mirror/pkg/utils"

	"google.golang.org/grpc"
)

// RegisterDestinationEndpoint registers destination
func RegisterDestinationEndpoint(host string, port uint16) (*grpc.ClientConn, error) {

	return utils.GetGRPCClient(conf.GRPCService{
		Host: &host,
		Port: port,
	})

}
