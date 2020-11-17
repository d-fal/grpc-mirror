package mirrors

import (
	"context"
	"fmt"
	"grpc-mirror/internal/broker"
	"grpc-mirror/pkg/libs"
	topupService "grpc-mirror/pkg/protobufs/destinationpb"
	pb "grpc-mirror/pkg/protobufs/mirrorpb"

	"github.com/logrusorgru/aurora"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// Mirror object
type Mirror struct {
	Endpoints GRPCEndpoint
	Broker    *broker.Broker
}

// SetupMirrors prepares a mirror
func SetupMirrors() *Mirror {
	endpoints := GRPCEndpoint{}
	return &Mirror{
		Endpoints: endpoints,
	}
}

// SetBroker sets a broker to be used in publish
func (m *Mirror) SetBroker(brkr *broker.Broker) *Mirror {
	m.Broker = brkr

	return m
}

// AddClient appends grpc endpoints
func (m *Mirror) AddClient(conn *grpc.ClientConn, f interface{}, contains *bool,
	filters ...DispatchOptions) *Mirror {

	m.Endpoints[&Connection{Conn: conn, Filters: filters, Contains: contains}] = f
	return m
}

// Start activates grpc endpoints to receive messages from broker
func (m *Mirror) Start() {

	for conn, f := range m.Endpoints {
		go func(conn *Connection, f interface{}) {

			msgChannel := m.Broker.Subscribe()

			for {
				msg, ok := (<-msgChannel).(*libs.ChannelBeacon)
				if ok {

					if ok, redirect := conn.Lookup(msg.Message.Type); ok {
						if ok {
							DeliverMessage(msg, f, *conn, redirect)
						}
					}
				}
			}

		}(conn, f)
	}
}

// DeliverMessage delivers message
func DeliverMessage(msg *libs.ChannelBeacon, f interface{}, conn Connection, redirect bool) {

	defer func() {
		if r := recover(); r != nil {
			// println("Channel is closed!")
		}
	}()
	switch f.(type) {

	case func(context.Context,
		*topupService.TopupRequestType,
		...grpc.CallOption) (*topupService.TopupResponseType, error):

		fn, ok := f.(func(context.Context,
			*topupService.TopupRequestType,
			...grpc.CallOption) (*topupService.TopupResponseType, error))

		if ok {
			if conn.Conn.GetState() != connectivity.Ready {
				fn = topupService.NewTopupServiceClient(conn.Conn).Topup // edit
			}

			response, err := fn(context.Background(), msg.Message)

			if err != nil {
				fmt.Println("Sending error: ", err, aurora.Red(response))
				return
			}
			if redirect {
				if msg.Wait() {
					msg.Response <- response
				}
			}
		} else {
			fmt.Println("Sending failed!", msg.Message.Type)
		}

	case func(context.Context,
		*pb.TopupRequestType,
		...grpc.CallOption) (*pb.TopupResponseType, error):

		fn, ok := f.(func(context.Context,
			*pb.TopupRequestType,
			...grpc.CallOption) (*pb.TopupResponseType, error))

		if ok {
			if conn.Conn.GetState() != connectivity.Ready {
				fn = pb.NewTopupServiceClient(conn.Conn).Topup // edit
			}
			_msg := UpliftTopupRequest(msg.Message)
			response, err := fn(context.Background(), _msg)

			if err != nil {
				fmt.Println("Sending error: ", err, aurora.Red(response))
				return
			}

			if redirect {
				if msg.Wait() && response != nil {

					msg.Response <- ConvertTopupResponse(response)
				}
			}
		} else {
			fmt.Println("Sending failed!", msg.Message.Type)
		}

	}
}
