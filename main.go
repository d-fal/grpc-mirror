package main

import (
	"context"
	"grpc-mirror/internal/broker"
	"grpc-mirror/internal/control"
	"grpc-mirror/internal/mirrors"
	"grpc-mirror/internal/server"
	"grpc-mirror/pkg/conf"
	"grpc-mirror/pkg/logger"
	topupService "grpc-mirror/pkg/protobufs/destinationpb"
	pb "grpc-mirror/pkg/protobufs/mirrorpb"
	"grpc-mirror/pkg/utils"

	"go.uber.org/zap"
)

// RPCMessage rpc message that we are willing to transfer to destination while making some copies of
type RPCMessage func()

var (
	zapLogger         *zap.Logger
	systemWideCancel  context.CancelFunc
	systemWideContext context.Context
	has               = true
	hasnt             = false
)

func init() {
	systemWideContext, systemWideCancel = context.WithCancel(context.Background())
	control.SignalListener(systemWideContext, systemWideCancel)

	logger.SetLogPath(conf.GetConfigObject().GetServerConfig().LogFile)

}

func main() {
	defer systemWideCancel()
	// Register grpc destination which origin wants to connect to

	// Register mirrors
	broker := broker.SetupBroker()
	go broker.Start()
	defer broker.Stop()

	// Load grpc server to accept calls from client, we call it origin.
	grpcServer := server.RegisterServer("0.0.0.0", 50051).
		SetupBroker(broker).
		SetTLS(false)

	mirror := mirrors.SetupMirrors().SetBroker(broker)

	url1 := "127.0.0.1"

	orderMS, _ := utils.GetGRPCClient(conf.GRPCService{
		Host: &url1,
		Port: 50052,
	})
	orderCtx := pb.NewTopupServiceClient(orderMS)

	loggerMS, _ := utils.GetGRPCClient(conf.GRPCService{
		Host: &url1,
		Port: 50056,
	})
	loggerCtx := topupService.NewTopupServiceClient(loggerMS)

	defer func() {
		orderMS.Close()
		loggerMS.Close()
	}()

	mirror.
		AddClient(orderMS, orderCtx.Topup, &has,
			mirrors.DispatchOptions{MessageType: 412, RedirectResponse: true}).
		AddClient(loggerMS, loggerCtx.Topup, &has,
			mirrors.DispatchOptions{MessageType: 412, RedirectResponse: true}).
		Start()

	if err := grpcServer.SetupServerConnection(topupService.RegisterTopupServiceServer); err != nil {
		return /* GRPC Services are not started */
	}

	// grpc server
	conf.SetLogPath(conf.GetConfigObject().GetServerConfig().LogFile)

	zapLogger = logger.GetZapLogger("incoming_messages")
	grpcServer.SetLogger(zapLogger)
	grpcServer.Listen(systemWideContext)
}
