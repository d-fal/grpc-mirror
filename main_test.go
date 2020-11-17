package main

import (
	"context"
	"grpc-mirror/internal/broker"
	"grpc-mirror/pkg/conf"
	topupService "grpc-mirror/pkg/protobufs/destinationpb"
	"grpc-mirror/pkg/utils"
	"sync"
	"testing"
	"time"

	"github.com/logrusorgru/aurora"
	"google.golang.org/grpc/connectivity"
)

func TestServer(t *testing.T) {
	t.Parallel()

	pool := &sync.Pool{}

	broker := broker.SetupBroker()
	go broker.Start()
	defer broker.Stop()

	url1 := "127.0.0.1"

	// grpcServer := server.RegisterServer("0.0.0.0", 50050).
	// 	SetupBroker(broker).
	// 	SetTLS(false)
	// waitC := make(chan struct{})

	// grpcServer.SetupServerConnection(destinationpb.RegisterTopupServiceServer)
	// go func() {
	// 	if err := grpcServer.Listen(systemWideContext); err != nil {
	// 		close(waitC)
	// 	}

	// }()
	mirror, _ := utils.GetGRPCClient(conf.GRPCService{
		Host: &url1,
		Port: 50051,
	})
	defer mirror.Close()

	for i := 0; i < 100; i++ {

		go func(i int) {
			if mirror.GetState() == connectivity.Ready {

				pool.Put(topupService.NewTopupServiceClient(mirror))

				req := &topupService.TopupRequestType{
					Username: "98092",
					Password: "1809209",
					Msisdn:   "1121",
					Pin:      "12111",
					MobNo:    "9120249937",
					Amount:   int32(i),
					Desc:     "No desc",
					AddData:  "Some data",
					Type:/*int32(410 + i%4)*/ 412,
				}
				if ctx, ok := pool.Get().(topupService.TopupServiceClient); ok {
					_, err := ctx.Topup(context.TODO(), req)
					if err != nil {
						t.Errorf("Error sending request %v\n", err)

					}
				}

			} else {
				t.Log("Channel is ", aurora.Yellow(mirror.GetState()), i)
			}
		}(i)
		time.Sleep(time.Millisecond * 1000)
	}

}
