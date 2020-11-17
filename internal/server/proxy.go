package server

import (
	"context"
	"grpc-mirror/pkg/libs"
	topupService "grpc-mirror/pkg/protobufs/destinationpb"
	"time"
)

/*
This is the service we have proxied, this service should alert others on changes occur
elsewehere in the app
*/

func publish(proxy *Proxy) {

}

var i int
var total int

// Topup topup unary service listener
func (srv *Server) Topup(ctx context.Context,
	req *topupService.TopupRequestType) (*topupService.TopupResponseType, error) {

	var response interface{}

	messageToDispatch := &libs.ChannelBeacon{
		Message:  req,
		Response: make(chan interface{}),
	}
	srv.logger.AddOne("request", req)

	srv.Broker.Publish(messageToDispatch)

	/* graceful stop implemented as well */
	select {

	case <-srv.gctx.Done():
		srv.grpcServer.GracefulStop()

	case response = <-messageToDispatch.Response:

		srv.logger.AddOne("response", response).Commit("success")

	case <-time.After(time.Second * 30): // open to ?

		srv.logger.AddOne("response", response).Commit("timeout")

	}

	messageToDispatch.Done()

	if _response, ok := response.(*topupService.TopupResponseType); ok {
		return _response, nil
	}

	return &topupService.TopupResponseType{
		RespCode: 1, // sample response
	}, nil

}
