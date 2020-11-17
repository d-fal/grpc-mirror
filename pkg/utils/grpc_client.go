package utils

import (
	"fmt"
	"grpc-mirror/pkg/conf"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GetGRPCClient get grpc client and
func GetGRPCClient(grpcConf conf.GRPCService) (*grpc.ClientConn, error) {

	var (
		connstr string
	)

	creds := []grpc.DialOption{}

	if grpcConf.TLSEnable {

		_creds, sslErr := credentials.NewClientTLSFromFile(*grpcConf.SSL.Ca, "")

		if sslErr != nil {
			log.Println("Error reading ssl ca ", sslErr)
			return nil, sslErr
		}

		creds = append(creds, grpc.WithTransportCredentials(_creds))

	} else {
		creds = append(creds, grpc.WithInsecure())
	}

	connstr = fmt.Sprintf("%s:%d", *grpcConf.Host, grpcConf.Port)

	conn, err := grpc.Dial(connstr, creds...)

	if err != nil {
		log.Printf("Could not connect: %v\n", err)
		return nil, err
	}
	return conn, nil
}
