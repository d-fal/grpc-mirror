package mirrors

import (
	"google.golang.org/grpc"
)

// GRPCEndpoint grpc endpoint
type GRPCEndpoint map[*Connection]interface{}

// Connection is the container that keeps connection
type Connection struct {
	Conn    *grpc.ClientConn
	Filters []DispatchOptions
	// Contains determines if a message with specified type should have a type or not.
	// This one is a little tricky. So here is the example:
	// {MessageType: 1, Contains: true}
	// The above example is equivalent to:
	// Select all the messages where MessageType=1
	// another example:
	// {MessageType: 3, Contains: false}
	// Select * messages where MessageType is not 3
	// To pass all the messages without any kinds of filters, just let it be nil
	Contains *bool
}

// DispatchOptions dispatch options
/*
									TOPUP (Recv types {1},{2,RedirectResponse}) TOPUP
								   /
 ______  Msg type: 3 ___________ /________ Mirror1 (Recv types {1},{2},{3}) KAFKA
|BPM   | --->-->--> |Dispatcher |/
|______|  ----<---<-|___________|\________ Mirror2 (Recv type {1} only) ELASTICSEARCH
		Rsp from Reciver2   	 \
								  \ Receiver2 (Recv types {2},{3,RedirectResponse}) BILL

*/
type DispatchOptions struct {
	MessageType int32

	// RedirectResponse is the hack to redirect responses
	// If Contains in Connection block was true, then it looks into this flag
	// If RedirectResponse was true, the dispatcher sends back te recipent
	// response to the client that has made the request
	RedirectResponse bool
}

// Lookup looks up into filterd types and prevents message dispatch if filter is specified
// This method searchs inside the filtered types and send back the result of search
// in accord with specified filter
// Assume we have a message with type:2 and a registered recipient with
// the following Dispatch options
// Option 1: {Type:1  RedirectResponse: false}
// Option 2: {Type:2  RedirectResponse: true}
// .....
// Since we have message type=2, Option 2 asserts that we should block the message that
// has this type. So, it will be blocked
// If a message was not specified in any of
func (c *Connection) Lookup(filter int32) (bool, bool) {
	if c.Contains != nil {
		for _, f := range c.Filters {
			if f.MessageType == filter {
				return *c.Contains, f.RedirectResponse
			}
		}
		return !*c.Contains, false
	}
	return true, false
}
